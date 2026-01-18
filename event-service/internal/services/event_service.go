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
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
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
		return nil, err
	}
	return event, nil
}

func (s *eventService) GetEvent(id uint) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, e.ErrEventNotFound
	}
	return event, nil
}

func (s *eventService) DeleteEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return e.ErrEventNotFound
	}

	if event.Status != string(dto.Draft) {
		return e.ErrEventIsNotDraft
	}
	return s.eventRepo.Delete(id)
}

func (s *eventService) UpdateEvent(req dto.UpdateEventRequest, id uint) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
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
		return nil, err
	}
	return event, nil
}

func (s *eventService) ListEvents(query dto.EventListQuery) ([]models.Event, error) {
	return s.eventRepo.List(query)
}

func (s *eventService) PublishEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return e.ErrEventNotFound
	}

	if event.Status != string(dto.Draft) {
		return e.ErrEventIsNotDraft
	}

	event.Status = string(dto.Published)

	return s.eventRepo.Update(event)
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
	return s.eventRepo.GetByUserID(userID)
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
