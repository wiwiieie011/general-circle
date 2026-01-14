package services

import (
	"event-service/internal/dto"
	"event-service/internal/models"
	"event-service/internal/repository"
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
}

type eventService struct {
	eventRepo    repository.EventRepository
	categoryRepo repository.CategoryRepository
}

func NewEventService(
	eventRepo repository.EventRepository,
	categoryRepo repository.CategoryRepository,
) EventService {
	return &eventService{eventRepo: eventRepo, categoryRepo: categoryRepo}
}

func (s *eventService) CreateEvent(req dto.CreateEventRequest) (*models.Event, error) {
	if req.Title == "" {
		return nil, dto.ErrEmptyTitle
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			return nil, dto.ErrCategoryNotFound
		}
	}

	req.Status = dto.Draft

	event := &models.Event{
		Title:      req.Title,
		Status:     string(req.Status),
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
		return nil, dto.ErrEventNotFound
	}
	return event, nil
}

func (s *eventService) DeleteEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return dto.ErrEventNotFound
	}

	if event.Status != string(dto.Draft) {
		return dto.ErrEventIsNotDraft
	}
	return s.eventRepo.Delete(id)
}

func (s *eventService) UpdateEvent(req dto.UpdateEventRequest, id uint) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, dto.ErrEventNotFound
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			return nil, dto.ErrEmptyTitle
		}
		event.Title = trimmed
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
		return dto.ErrEventNotFound
	}

	event.Status = string(dto.Published)

	return s.eventRepo.Update(event)
}

func (s *eventService) CancelEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return dto.ErrEventNotFound
	}

	if event.Status != string(dto.Published) {
		return dto.ErrEventIsNotPublished
	}

	event.Status = string(dto.Cancelled)

	return s.eventRepo.Update(event)
}
