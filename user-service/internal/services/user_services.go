package services

import (
	"errors"

	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(email, password, firstName, lastName string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	BecomeOrganizer(userID uint) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(email, password, firstName, lastName string) (*models.User, error) {
_, err := s.userRepo.GetByEmail(email)
if err == nil {
	return nil, errors.New("пользователь с таким email уже существует")
}

if !errors.Is(err, gorm.ErrRecordNotFound) {
	return nil, err
}


	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("неверный email или пароль")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) BecomeOrganizer(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Role == models.RoleOrganizer {
		return nil, errors.New("пользователь уже организатор")
	}

	user.Role = models.RoleOrganizer

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
