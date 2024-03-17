package userhandler

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"tiny/internal/models"
)

type JWTManager struct {
	secret        string
	tokenDuration time.Duration
}

type UseClaims struct {
	jwt.RegisteredClaims
	Login string
	Id    int
}

func NewJWTManager(secret string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secret: secret, tokenDuration: tokenDuration}
}

func (m *JWTManager) Generate(user models.User) (string, error) {
	claims := UseClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
		},
		Login: user.Login,
		Id:    user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTManager) Parse(Token string) (*UseClaims, error) {
	token, err := jwt.ParseWithClaims(
		Token,
		&UseClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token string method")
			}

			return []byte(m.secret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UseClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
