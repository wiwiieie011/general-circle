package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token expired")
	ErrUnexpectedSign = errors.New("unexpected signing method")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	UserID uint      `json:"user_id"`
	Role   string    `json:"role,omitempty"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey  string
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewTokenManager(
	secretKey string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	issuer string,
) *TokenManager {
	return &TokenManager{
		secretKey:  secretKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
}


func (tm *TokenManager) GenerateAccessToken(
	userID uint,
	role string,
) (string, error) {
	return tm.generateToken(userID, role, AccessToken, tm.accessTTL)
}

func (tm *TokenManager) GenerateRefreshToken(
	userID uint,
) (string, error) {
	return tm.generateToken(userID, "", RefreshToken, tm.refreshTTL)
}

func (tm *TokenManager) generateToken(
	userID uint,
	role string,
	tokenType TokenType,
	ttl time.Duration,
) (string, error) {

	now := time.Now()

	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}


func (tm *TokenManager) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSign
			}
			return []byte(tm.secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
