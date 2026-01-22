package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	e "user-service/internal/errors"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokenManager *utils.TokenManager
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenManager *utils.TokenManager,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}


func (s *AuthService) Register(email string,password string,firstName string,lastName string,) (*models.User, string, string, error) {

	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, "", "", e.ErrEmailAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, "", "", err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         models.RoleUser,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", "", err
	}

	accessToken, err :=
		s.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err :=
		s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}


func (s *AuthService) Login(email string,password string,) (*models.User, string, string, error) {

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", "", e.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", "", e.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, "", "", e.ErrInvalidCredentials
	}

	accessToken, err :=
		s.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err :=
		s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string,) (string, string, error) {

	claims, err := s.tokenManager.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.Type != utils.RefreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", e.ErrUserNotFound
	}

	if !user.IsActive {
		return "", "", e.ErrUserInactive
	}

	newAccessToken, err :=
		s.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err :=
		s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

	func (s *AuthService) IssueAccessTokenForUser(	user *models.User,) (string, error) {
		return s.tokenManager.GenerateAccessToken(
			user.ID,
			string(user.Role),
		)
	}

