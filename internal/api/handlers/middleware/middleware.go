package middleware

import (
	"time"
	"tiny/internal/usecase"
	userhandler "tiny/pkg"
)

type Middleware struct {
	// TODO: I change session type
	users      usecase.Users
	sessions   usecase.Sessions
	jwt        *userhandler.JWTManager
	sessionTTL time.Duration
	tokenTTL   time.Duration
}

func NewMiddleware(users usecase.Users, sessions usecase.Sessions, jwtManager userhandler.JWTManager, tokenTTL, sessionTTL time.Duration) *Middleware {
	return &Middleware{
		users:      users,
		sessions:   sessions,
		jwt:        &jwtManager,
		sessionTTL: sessionTTL,
		tokenTTL:   tokenTTL,
	}
}
