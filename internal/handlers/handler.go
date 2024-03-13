package handlers

import (
	actorshandler "tiny/internal/handlers/actors"
	filmshandler "tiny/internal/handlers/films"
	userhandler "tiny/internal/handlers/users"
)

type Handler struct {
	actors actorshandler.ActorsHandler
	films  filmshandler.FilmsHandler
	users  userhandler.UserHandler
}

func NewHandler(a actorshandler.ActorsHandler, f filmshandler.FilmsHandler, u userhandler.UserHandler) *Handler {
	return &Handler{a, f, u}
}
