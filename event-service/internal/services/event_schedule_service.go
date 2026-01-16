package services

import (
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/models"
	"event-service/internal/repository"
)

type EventScheduleService interface {
	GetScheduleByEventID(eventID uint) ([]models.EventSchedule, error)
	CreateScheduleForEvent(eventID uint, req dto.CreateScheduleRequest) (*models.EventSchedule, error)
}

type eventScheduleService struct {
	eventScheduleRepo repository.EventScheduleRepository
	eventRepo         repository.EventRepository
}

func NewEventScheduleService(
	eventScheduleRepo repository.EventScheduleRepository,
	eventRepo repository.EventRepository,
) EventScheduleService {
	return &eventScheduleService{
		eventScheduleRepo: eventScheduleRepo,
		eventRepo:         eventRepo,
	}
}

func (s *eventScheduleService) GetScheduleByEventID(eventID uint) ([]models.EventSchedule, error) {
	if _, err := s.eventRepo.GetByID(eventID); err != nil {
		return nil, e.ErrEventNotFound
	}

	schedules, err := s.eventScheduleRepo.GetByEventID(eventID)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *eventScheduleService) CreateScheduleForEvent(
	eventID uint,
	req dto.CreateScheduleRequest,
) (*models.EventSchedule, error) {
	if _, err := s.eventRepo.GetByID(eventID); err != nil {
		return nil, e.ErrEventNotFound
	}
	if req.ActivityName == "" {
		return nil, e.ErrEmptyActivityName
	}

	if req.Speaker == "" {
		return nil, e.ErrEmptySpeaker
	}

	if !req.StartAt.Before(req.EndAt) {
		return nil, e.ErrNotCorrectScheduleTime
	}

	schedule := &models.EventSchedule{
		EventID:      eventID,
		ActivityName: req.ActivityName,
		Speaker:      req.Speaker,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
	}

	if err := s.eventScheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}
