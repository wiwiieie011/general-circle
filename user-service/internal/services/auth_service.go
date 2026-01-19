package services

import (
	"errors"
	"time"

	"user-service/internal/config"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("invalid user inactive")
)

type AuthService struct {
	repo repository.UserRepository
	jwt  config.JWTConfig
}

func NewAuthService(repo repository.UserRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwtCfg,
	}
}

func (s *AuthService) Login(email, password string) (string, string, *models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", "", nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	accessTTL := time.Duration(s.jwt.AccessTTL) * time.Minute
	refreshTTL := time.Duration(s.jwt.RefreshTTL) * 24 * time.Hour

	access, err := utils.GenerateToken(
		user.ID,
		string(user.Role),
		s.jwt.AccessSecret,
		accessTTL,
	)
	if err != nil {
		return "", "", nil, err
	}

	refresh, err := utils.GenerateToken(
		user.ID,
		string(user.Role),
		s.jwt.RefreshSecret,
		refreshTTL,
	)
	if err != nil {
		return "", "", nil, err
	}

	return access, refresh, user, nil
}
