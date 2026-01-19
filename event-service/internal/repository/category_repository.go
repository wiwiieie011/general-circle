package repository

import (
	"errors"
	"event-service/internal/models"

	e "event-service/internal/errors"
	"log/slog"

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
	db     *gorm.DB
	logger *slog.Logger
}

func NewCategoryRepository(db *gorm.DB, logger *slog.Logger) CategoryRepository {
	return &gormCategoryRepository{db: db, logger: logger}
}

func (r *gormCategoryRepository) Create(category *models.Category) error {
	if category == nil {
		return e.ErrCategoryIsNil
	}
	r.logger.Debug("creating category", slog.String("name", category.Name))
	if err := r.db.Create(category).Error; err != nil {
		r.logger.Error("failed to create category", "error", err)
		return err
	}
	return nil
}

func (r *gormCategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category

	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("category not found by id", slog.Int("id", int(id)))
			return nil, e.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) Delete(id uint) error {
	r.logger.Debug("deleting category", slog.Int("id", int(id)))
	if err := r.db.Delete(&models.Category{}, id).Error; err != nil {
		r.logger.Error("failed to delete category", "error", err, "id", id)
		return err
	}
	return nil
}

func (r *gormCategoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category

	if err := r.db.Where("name = ?", name).
		First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("category not found by name", slog.String("name", name))
			return nil, e.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category

	if err := r.db.Find(&categories).Error; err != nil {
		r.logger.Error("failed to list categories", "error", err)
		return nil, err
	}
	return categories, nil
}
