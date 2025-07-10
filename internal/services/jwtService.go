package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	jwtSecret []byte
}

func NewJwtService(jwtSecret string) *JwtService {
	return &JwtService{jwtSecret: []byte(jwtSecret)}
}

func (s *JwtService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	return true, nil
}
