package actorshandler

import "tiny/internal/usecase"

type ActorsHandler struct {
	actors usecase.Actors
}

func NewActorsHandler(a usecase.Actors) *ActorsHandler {
	return &ActorsHandler{a}
}
