package repository

import (
	"event-service/internal/dto"
	"event-service/internal/models"
	"strings"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	GetByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
	List(query dto.EventListQuery) ([]models.Event, error)
}

type gormEventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &gormEventRepository{db: db}
}

func (r *gormEventRepository) Create(event *models.Event) error {
	if event == nil {
		return dto.ErrEventIsNil
	}

	return r.db.Create(event).Error
}

func (r *gormEventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event

	if err := r.db.Preload("Schedule").
		First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *gormEventRepository) Update(event *models.Event) error {
	if event == nil {
		return dto.ErrEventIsNil
	}

	return r.db.Save(event).Error
}

func (r *gormEventRepository) Delete(id uint) error {
	return r.db.Preload("Schedule").Delete(&models.Event{}, id).Error
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

	if err := db.Preload("Schedule").
		Order(sortField + " " + order).
		Limit(query.Limit).
		Offset(offset).
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
