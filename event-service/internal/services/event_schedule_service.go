package services

import (
	"event-service/internal/dto"
	"event-service/internal/models"
	"event-service/internal/repository"
)

type EventScheduleService interface {
	CreateSchedule(req dto.CreateScheduleRequest) (*models.EventSchedule, error)
	GetSchedule(id uint) (*models.EventSchedule, error)
}

type eventScheduleService struct {
	eventScheduleRepo repository.EventScheduleRepository
}

func NewEventScheduleService(
	eventScheduleRepo repository.EventScheduleRepository,
) EventScheduleService {
	return &eventScheduleService{eventScheduleRepo: eventScheduleRepo}
}

func (s *eventScheduleService) CreateSchedule(req dto.CreateScheduleRequest) (*models.EventSchedule, error) {
	if req.ActivityName == "" {
		return nil, dto.ErrEmptyActivityName
	}

	if req.EventID < 1 {
		return nil, dto.ErrNotCorrectID
	}

	if req.Speaker == "" {
		return nil, dto.ErrEmptySpeaker
	}

	if !req.StartAt.Before(req.EndAt) {
		return nil, dto.ErrNotCorrectScheduleTime 
	}

	schedule := &models.EventSchedule{
		EventID:      req.EventID,
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

func (s *eventScheduleService) GetSchedule(id uint) (*models.EventSchedule, error) {
	schedule, err := s.eventScheduleRepo.GetByID(id)
	if err != nil {
		return nil, dto.ErrEventScheduleNotFound
	}
	return schedule, nil
}
