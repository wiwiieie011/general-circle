package repository

import (
	"event-service/internal/dto"
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
		return dto.ErrEventScheduleIsNil
	}
	return r.db.Create(schedule).Error
}

func (r *gormScheduleRepository) GetByID(id uint) (*models.EventSchedule, error) {
	var schedule models.EventSchedule

	if err := r.db.First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *gormScheduleRepository) GetByEventID(eventID uint) ([]models.EventSchedule, error) {
	var schedules []models.EventSchedule

	if err := r.db.Where("event_id = ?", eventID).
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
