package repository

import (
	e "event-service/internal/errors"
	"event-service/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type EventScheduleRepository interface {
	Create(schedule *models.EventSchedule) error
	GetByID(id uint) (*models.EventSchedule, error)
	GetByEventID(eventID uint) ([]models.EventSchedule, error)
}

type gormScheduleRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewEventScheduleRepository(db *gorm.DB, logger *slog.Logger) EventScheduleRepository {
	return &gormScheduleRepository{db: db, logger: logger}
}

func (r *gormScheduleRepository) Create(schedule *models.EventSchedule) error {
	if schedule == nil {
		return e.ErrEventScheduleIsNil
	}
	r.logger.Debug("creating schedule", slog.Int("event_id", int(schedule.EventID)))
	if err := r.db.Create(schedule).Error; err != nil {
		r.logger.Error("failed to create schedule", "error", err)
		return err
	}
	return nil
}

func (r *gormScheduleRepository) GetByID(id uint) (*models.EventSchedule, error) {
	var schedule models.EventSchedule

	if err := r.db.First(&schedule, id).Error; err != nil {
		r.logger.Error("failed to get schedule by id", "error", err, "id", id)
		return nil, err
	}
	return &schedule, nil
}

func (r *gormScheduleRepository) GetByEventID(eventID uint) ([]models.EventSchedule, error) {
	var schedules []models.EventSchedule

	if err := r.db.Where("event_id = ?", eventID).
		Find(&schedules).Error; err != nil {
		r.logger.Error("failed to get schedules by event", "error", err, "event_id", eventID)
		return nil, err
	}
	return schedules, nil
}
