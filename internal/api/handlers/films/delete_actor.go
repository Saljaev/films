package filmshandler

import (
	"context"
	"net/http"
	"tiny/internal/api/utilapi"
)

type FilmDeleteActorResponse struct {
	FilmID  int `json:"film_id"`
	ActorID int `json:"actor_id"`
}

type FilmDeleteActorRequest struct {
	FilmID  int `json:"film_id"`
	ActorID int `json:"actor_id"`
}

func (req *FilmDeleteActorRequest) IsValid() bool {
	return req.ActorID > 0 && req.FilmID > 0
}

func (h *FilmsHandler) DeleteActor(ctx *utilapi.APIContext) {
	var req FilmDeleteActorRequest

	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	ctx.Info("request decoded", "request", req)

	film, err := h.films.GetById(context.Background(), req.FilmID)
	if err != nil {
		ctx.Error("failed to delete film by id", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid film id")
		return
	}

	err = h.films.DeleteActor(context.Background(), film.Id, req.ActorID)
	if err != nil {
		ctx.Error("failed to delete actor from film", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film deleted", "film_id", req.FilmID)
	ctx.SuccessWithData(FilmDeleteActorResponse{FilmID: req.FilmID, ActorID: req.ActorID})
}
