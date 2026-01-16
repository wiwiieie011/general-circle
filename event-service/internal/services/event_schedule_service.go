package services

import (
	"event-service/internal/dto"
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
	return &eventScheduleService{eventScheduleRepo: eventScheduleRepo}
}

func (s *eventScheduleService) GetScheduleByEventID(eventID uint) ([]models.EventSchedule, error) {
	if _, err := s.eventRepo.GetByID(eventID); err != nil {
		return nil, dto.ErrEventNotFound
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
		return nil, dto.ErrEventNotFound
	}
	if req.ActivityName == "" {
		return nil, dto.ErrEmptyActivityName
	}

	if req.Speaker == "" {
		return nil, dto.ErrEmptySpeaker
	}

	if !req.StartAt.Before(req.EndAt) {
		return nil, dto.ErrNotCorrectScheduleTime
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
