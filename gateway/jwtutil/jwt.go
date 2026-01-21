package jwtutil

import (
	"os"
	"github.com/golang-jwt/jwt/v5"
)

func getSecret() []byte {
	return []byte(os.Getenv("SUPER_SECRET_KEY"))
}

type Claims struct {
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}


func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return getSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
