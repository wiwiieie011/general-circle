package repository

import (
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	GetByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
	List(query dto.EventListQuery) ([]models.Event, error)
	GetByUserID(userID uint) ([]models.Event, error)
	GetEventStartingTomorrow() ([]models.Event, error)
}

type gormEventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &gormEventRepository{db: db}
}

func (r *gormEventRepository) Create(event *models.Event) error {
	if event == nil {
		return e.ErrEventIsNil
	}
	return r.db.Create(event).Error
}

func (r *gormEventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event

	if err := r.db.Preload("Category").
		Preload("Schedule").
		First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *gormEventRepository) Update(event *models.Event) error {
	if event == nil {
		return e.ErrEventIsNil
	}
	return r.db.Save(event).Error
}

func (r *gormEventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

func (r *gormEventRepository) List(query dto.EventListQuery) ([]models.Event, error) {
	db := r.db.Model(&models.Event{})

	if query.Title != "" {
		db = db.Where("title ILIKE ?", "%"+query.Title+"%")
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(query.SortOrder))

	validSortFields := map[string]string{
		"title":      "title",
		"created_at": "created_at",
	}

	sortField, ok := validSortFields[sortBy]
	if !ok {
		sortField = "created_at"
	}

	validOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	order, ok := validOrders[sortOrder]
	if !ok {
		order = "DESC"
	}

	if query.Page < 1 {
		query.Page = dto.DefaultPage
	}

	if query.Limit < 1 {
		query.Limit = dto.DefaultLimit
	}

	offset := (query.Page - 1) * query.Limit

	var events []models.Event

	if err := db.Preload("Category").
		Preload("Schedule").
		Order(sortField + " " + order).
		Limit(query.Limit).
		Offset(offset).
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *gormEventRepository) GetByUserID(userID uint) ([]models.Event, error) {
	var events []models.Event

	if err := r.db.Where("user_id = ?", userID).
		Preload("Category").
		Preload("Schedule").
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *gormEventRepository) GetEventStartingTomorrow() ([]models.Event, error) {
	var events []models.Event

	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrowStart := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	tomorrowEnd := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 999999999, tomorrow.Location())

	var eventsIDs []uint
	err := r.db.Model(&models.EventSchedule{}).
		Select("DISTINCT event_id").
		Where("start_at >= ? AND start_at <= ?", tomorrowStart, tomorrowEnd).
		Pluck("event_id", &eventsIDs).Error
	if err != nil {
		return nil, err
	}

	if len(eventsIDs) == 0 {
		return []models.Event{}, nil
	}

	err = r.db.Where("status = ?", "published").
		Where("id IN ?").
		Preload("Category").
		Preload("Schedule").
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}
