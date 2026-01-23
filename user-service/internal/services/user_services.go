package services

import (
	e "user-service/internal/errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type UserService interface {
	GetByID(id uint) (*models.User, error)
	UpdateProfile(id uint, firstName, lastName string) (*models.User, error)
	BecomeOrganizer(id uint) (*models.User, error)
	GetByIDs(id uint) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, e.ErrUserNotFound
	}

	if user.Role != models.RoleOrganizer {
		return nil, e.ErrNotOrganizer
	}

	if !user.IsActive {
		return nil, e.ErrUserInactive
	}

	return user, nil
}
func (s *userService) GetByIDs(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, e.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, e.ErrUserInactive
	}

	return user, nil
}

func (s *userService) UpdateProfile(
	id uint,
	firstName string,
	lastName string,
) (*models.User, error) {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, e.ErrUserNotFound
	}

	user.FirstName = firstName
	user.LastName = lastName

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) BecomeOrganizer(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, e.ErrUserNotFound
	}

	if user.Role == models.RoleOrganizer {
		return nil, e.ErrAlreadyOrganizer
	}

	user.Role = models.RoleOrganizer
	
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
