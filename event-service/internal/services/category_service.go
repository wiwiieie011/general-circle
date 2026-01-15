package services

import (
	"errors"
	"event-service/internal/dto"
	"event-service/internal/models"
	"event-service/internal/repository"
)

type CategoryService interface {
	CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error)
	GetCategory(id uint) (*models.Category, error)
	DeleteCategory(id uint) error
	ListCategories() ([]models.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error) {
	if req.Name == "" {
		return nil, dto.ErrEmptyName
	}

	existing, err := s.categoryRepo.GetByName(req.Name)
	if err != nil {
		if !errors.Is(err, dto.ErrCategoryNotFound) {
			return nil, err
		}
	} else if existing != nil {
		return nil, dto.ErrCategoryNameExists
	}

	category := &models.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetCategory(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, dto.ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	if _, err := s.categoryRepo.GetByID(id); err != nil {
		return dto.ErrCategoryNotFound
	}
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) ListCategories() ([]models.Category, error) {
	return s.categoryRepo.List()
}
