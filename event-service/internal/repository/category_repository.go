package repository

import (
	"errors"
	"event-service/internal/models"

	e "event-service/internal/errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	Delete(id uint) error
	GetByName(name string) (*models.Category, error)
	List() ([]models.Category, error)
}

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) Create(category *models.Category) error {
	if category == nil {
		return e.ErrCategoryIsNil
	}
	return r.db.Create(category).Error
}

func (r *gormCategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category

	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *gormCategoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category

	if err := r.db.Where("name = ?", name).
		First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
