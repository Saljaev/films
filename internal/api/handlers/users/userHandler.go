package userhandler

import (
	"time"
	"tiny/internal/usecase"
	"tiny/pkg"
)

type UserHandler struct {
	users      usecase.Users
	sessions   usecase.Sessions
	jwt        pkg.JWTManager
	tokenTTL   time.Duration
	sessionTTL time.Duration
}

func NewUserHandler(u usecase.Users, s usecase.Sessions, j pkg.JWTManager, tTTL, sTTL time.Duration) *UserHandler {
	return &UserHandler{u, s, j, tTTL, sTTL}
}
