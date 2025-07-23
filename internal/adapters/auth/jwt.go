package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"starter/config"
	"starter/internal/core/auth"
	"time"
)

type Generator struct {
}

func NewTokenGenerator() *Generator {
	return &Generator{}
}

func (generator *Generator) Generate(user auth.AuthenticatedUser) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
