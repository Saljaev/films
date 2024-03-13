package userhandler

import (
	"time"
	"tiny/internal/usecase"
)

type UserHandler struct {
	users      usecase.Users
	sessions   usecase.Sessions
	jwt        JWTManager
	tokenTTL   time.Duration
	sessionTTL time.Duration
}

func NewUserHandler(u usecase.Users, s usecase.Sessions, j JWTManager, tTTL, sTTL time.Duration) *UserHandler {
	return &UserHandler{u, s, j, tTTL, sTTL}
}
