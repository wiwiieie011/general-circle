package services

import (
	"context"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/kafka"
	"event-service/internal/models"
	"event-service/internal/repository"
	"log/slog"
	"strings"
)

type EventService interface {
	CreateEvent(req dto.CreateEventRequest) (*models.Event, error)
	GetEvent(id uint) (*models.Event, error)
	DeleteEvent(id uint) error
	UpdateEvent(req dto.UpdateEventRequest, id uint) (*models.Event, error)
	ListEvents(query dto.EventListQuery) ([]models.Event, error)
	PublishEvent(id uint) error
	CancelEvent(id uint) error
	GetEventsByUserID(userID uint) ([]models.Event, error)
	SendEventReminders(ctx context.Context) error
}

type eventService struct {
	eventRepo     repository.EventRepository
	categoryRepo  repository.CategoryRepository
	kafkaProducer kafka.EventProducer
	logger        *slog.Logger
}

func NewEventService(
	eventRepo repository.EventRepository,
	categoryRepo repository.CategoryRepository,
	kafkaProducer kafka.EventProducer,
	logger *slog.Logger,
) EventService {
	return &eventService{
		eventRepo:     eventRepo,
		categoryRepo:  categoryRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (s *eventService) CreateEvent(req dto.CreateEventRequest) (*models.Event, error) {
	s.logger.Debug("CreateEvent called",
		slog.String("title", req.Title),
		slog.Int("user_id", int(req.UserID)),
	)
	if req.CategoryID != nil {
		s.logger.Debug("CreateEvent has category", slog.Int("category_id", int(*req.CategoryID)))
	}
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			s.logger.Warn("category not found on create event", "category_id", *req.CategoryID)
			return nil, e.ErrCategoryNotFound
		}
	}

	event := &models.Event{
		Title:      strings.TrimSpace(req.Title),
		Status:     string(dto.Draft),
		UserID:     req.UserID,
		Seats:      req.Seats,
		CategoryID: req.CategoryID,
	}

	if err := s.eventRepo.Create(event); err != nil {
		s.logger.Error("failed to create event", "error", err, "title", event.Title)
		return nil, err
	}
	s.logger.Info("event created", slog.Int("id", int(event.ID)), slog.String("title", event.Title))
	return event, nil
}

func (s *eventService) GetEvent(id uint) (*models.Event, error) {
	s.logger.Debug("GetEvent called", slog.Int("id", int(id)))
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("event not found", "id", id)
		return nil, e.ErrEventNotFound
	}
	s.logger.Debug("event loaded", slog.Int("id", int(event.ID)), slog.String("title", event.Title))
	return event, nil
}

func (s *eventService) DeleteEvent(id uint) error {
	s.logger.Debug("DeleteEvent called", slog.Int("id", int(id)))
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("event not found for delete", "id", id)
		return e.ErrEventNotFound
	}

	if event.Status != string(dto.Draft) {
		s.logger.Warn("attempt to delete non-draft event", "id", id, "status", event.Status)
		return e.ErrEventIsNotDraft
	}

	if err := s.eventRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete event", "error", err, "id", id)
		return err
	}
	s.logger.Info("event deleted", slog.Int("id", int(id)))
	return nil
}

func (s *eventService) UpdateEvent(req dto.UpdateEventRequest, id uint) (*models.Event, error) {
	s.logger.Debug("UpdateEvent called", slog.Int("id", int(id)))
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("event not found for update", "id", id)
		return nil, e.ErrEventNotFound
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			return nil, e.ErrEmptyTitle
		}
		event.Title = trimmed
	}

	if req.CategoryID != nil {
		if _, err := s.categoryRepo.GetByID(*req.CategoryID); err != nil {
			return nil, e.ErrCategoryNotFound
		}
		event.CategoryID = req.CategoryID
	}

	if req.Seats != nil {
		if *req.Seats < 1 {
			return nil, e.ErrNotCorrectNum
		}
		event.Seats = req.Seats
	}

	if req.UserID != nil {
		event.UserID = *req.UserID
	}

	if err := s.eventRepo.Update(event); err != nil {
		s.logger.Error("failed to update event", "error", err, "id", event.ID)
		return nil, err
	}
	s.logger.Info("event updated", slog.Int("id", int(event.ID)), slog.String("title", event.Title))
	return event, nil
}

func (s *eventService) ListEvents(query dto.EventListQuery) ([]models.Event, error) {
	s.logger.Debug("ListEvents called", slog.String("title", query.Title), slog.String("status", query.Status))
	events, err := s.eventRepo.List(query)
	if err != nil {
		s.logger.Error("failed to list events", "error", err)
		return nil, err
	}
	s.logger.Debug("ListEvents result", slog.Int("count", len(events)))
	return events, nil
}

func (s *eventService) PublishEvent(id uint) error {
	s.logger.Debug("PublishEvent called", slog.Int("id", int(id)))
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("event not found for publish", "id", id)
		return e.ErrEventNotFound
	}

	if event.Status != string(dto.Draft) {
		s.logger.Warn("attempt to publish non-draft event", "id", id, "status", event.Status)
		return e.ErrEventIsNotDraft
	}

	event.Status = string(dto.Published)

	if err := s.eventRepo.Update(event); err != nil {
		s.logger.Error("failed to publish event", "error", err, "id", id)
		return err
	}
	s.logger.Info("event published", slog.Int("id", int(id)))
	return nil
}

func (s *eventService) CancelEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return e.ErrEventNotFound
	}

	if event.Status != string(dto.Published) {
		return e.ErrEventIsNotPublished
	}

	event.Status = string(dto.Cancelled)

	if err := s.eventRepo.Update(event); err != nil {
		return err
	}

	ctx := context.Background()
	if err := s.kafkaProducer.SendEventCancelled(ctx, id); err != nil {
		s.logger.Error("failed to send event cancelled to kafka",
			"error", err,
			"event_id", id)
	}

	return nil
}

func (s *eventService) GetEventsByUserID(userID uint) ([]models.Event, error) {
	s.logger.Debug("GetEventsByUserID called", slog.Int("user_id", int(userID)))
	events, err := s.eventRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("failed to get events by user", "error", err, "user_id", userID)
		return nil, err
	}
	s.logger.Debug("GetEventsByUserID result", slog.Int("count", len(events)), slog.Int("user_id", int(userID)))
	return events, nil
}

func (s *eventService) SendEventReminders(ctx context.Context) error {
	events, err := s.eventRepo.GetEventStartingTomorrow()
	if err != nil {
		return err
	}

	for _, event := range events {
		if len(event.Schedule) == 0 {
			s.logger.Warn("event has no schedule", "event_id", event.ID)
			continue
		}

		firstActivity := event.Schedule[0]
		for _, schedule := range event.Schedule {
			if schedule.StartAt.Before(firstActivity.StartAt) {
				firstActivity = schedule
			}
		}

		if err := s.kafkaProducer.SendEventReminder(ctx, event.ID, event.Title, firstActivity.StartAt); err != nil {
			s.logger.Error("failed to send event reminder",
				"error", err,
				"event_id", event.ID)
		}
	}
	return nil
}
