package repository

import (
	e "event-service/internal/errors"
	"event-service/internal/models"

	"gorm.io/gorm"
)

type EventScheduleRepository interface {
	Create(schedule *models.EventSchedule) error
	GetByID(id uint) (*models.EventSchedule, error)
	GetByEventID(eventID uint) ([]models.EventSchedule, error)
}

type gormScheduleRepository struct {
	db *gorm.DB
}

func NewEventScheduleRepository(db *gorm.DB) EventScheduleRepository {
	return &gormScheduleRepository{db: db}
}

func (r *gormScheduleRepository) Create(schedule *models.EventSchedule) error {
	if schedule == nil {
		return e.ErrEventScheduleIsNil
	}
	return r.db.Create(schedule).Error
}

func (r *gormScheduleRepository) GetByID(id uint) (*models.EventSchedule, error) {
	var schedule models.EventSchedule

	if err := r.db.Preload("Event").First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *gormScheduleRepository) GetByEventID(eventID uint) ([]models.EventSchedule, error) {
	var schedules []models.EventSchedule

	if err := r.db.Where("event_id = ?", eventID).
		Preload("Event").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
