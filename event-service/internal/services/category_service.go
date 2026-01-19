package services

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/models"
	"event-service/internal/repository"
	"log/slog"
)

type CategoryService interface {
	CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error)
	GetCategory(id uint) (*models.Category, error)
	DeleteCategory(id uint) error
	ListCategories() ([]models.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	logger       *slog.Logger
}

func NewCategoryService(categoryRepo repository.CategoryRepository, logger *slog.Logger) CategoryService {
	return &categoryService{categoryRepo: categoryRepo, logger: logger}
}

func (s *categoryService) CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error) {
	s.logger.Debug("CreateCategory called", slog.String("name", req.Name))

	existing, err := s.categoryRepo.GetByName(req.Name)
	if err != nil {
		if !errors.Is(err, e.ErrCategoryNotFound) {
			return nil, err
		}
	} else if existing != nil {
		return nil, e.ErrCategoryNameExists
	}

	category := &models.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("failed to create category", "error", err, "name", req.Name)
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetCategory(id uint) (*models.Category, error) {
	s.logger.Debug("GetCategory called", slog.Int("id", int(id)))
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("category not found", "id", id)
		return nil, e.ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	s.logger.Debug("DeleteCategory called", slog.Int("id", int(id)))
	if _, err := s.categoryRepo.GetByID(id); err != nil {
		return e.ErrCategoryNotFound
	}
	if err := s.categoryRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete category", "error", err, "id", id)
		return err
	}
	return nil
}

func (s *categoryService) ListCategories() ([]models.Category, error) {
	return s.categoryRepo.List()
}