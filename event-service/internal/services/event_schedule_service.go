package services

import (
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/models"
	"event-service/internal/repository"
	"log/slog"
)

type EventScheduleService interface {
	GetScheduleByEventID(eventID uint) ([]models.EventSchedule, error)
	CreateScheduleForEvent(eventID uint, req dto.CreateScheduleRequest) (*models.EventSchedule, error)
}

type eventScheduleService struct {
	eventScheduleRepo repository.EventScheduleRepository
	eventRepo         repository.EventRepository
	logger            *slog.Logger
}

func NewEventScheduleService(
	eventScheduleRepo repository.EventScheduleRepository,
	eventRepo repository.EventRepository,
	logger *slog.Logger,
) EventScheduleService {
	return &eventScheduleService{
		eventScheduleRepo: eventScheduleRepo,
		eventRepo:         eventRepo,
		logger:            logger,
	}
}

func (s *eventScheduleService) GetScheduleByEventID(eventID uint) ([]models.EventSchedule, error) {
	s.logger.Debug("GetScheduleByEventID called", slog.Int("event_id", int(eventID)))
	if _, err := s.eventRepo.GetByID(eventID); err != nil {
		s.logger.Warn("event not found for schedule", "event_id", eventID)
		return nil, e.ErrEventNotFound
	}

	schedules, err := s.eventScheduleRepo.GetByEventID(eventID)
	if err != nil {
		s.logger.Error("failed to get schedules", "error", err, "event_id", eventID)
		return nil, err
	}
	return schedules, nil
}

func (s *eventScheduleService) CreateScheduleForEvent(
	eventID uint,
	req dto.CreateScheduleRequest,
) (*models.EventSchedule, error) {
	s.logger.Debug("CreateScheduleForEvent called", slog.Int("event_id", int(eventID)), slog.String("activity", req.ActivityName))
	if _, err := s.eventRepo.GetByID(eventID); err != nil {
		s.logger.Warn("event not found when creating schedule", "event_id", eventID)
		return nil, e.ErrEventNotFound
	}

	if !req.StartAt.Before(req.EndAt) {
		s.logger.Warn("invalid schedule time", "event_id", eventID)
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
		s.logger.Error("failed to create schedule", "error", err, "event_id", eventID)
		return nil, err
	}

	return schedule, nil
}
