package filmshandler

import (
	"context"
	"net/http"
	"strconv"
	"tiny/internal/api/utilapi"
)

type FilmDeleteResponse struct {
	FilmID int `json:"film_id"`
}

func (h *FilmsHandler) Delete(ctx *utilapi.APIContext) {
	rawFilmID := ctx.GetURLParam("id")

	filmID, err := strconv.Atoi(rawFilmID)
	if err != nil || filmID <= 0 {
		ctx.Error("invalid film id from url param", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid film id")
		return
	}

	ctx.Info("request decoded", "id", filmID)

	film, err := h.films.GetById(context.Background(), filmID)
	if err != nil {
		ctx.Error("failed to delete film by id", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid film id")
		return
	}

	err = h.films.Delete(context.Background(), film.Id)
	if err != nil {
		ctx.Error("failed to delete film", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film deleted", "film_id", filmID)
	ctx.SuccessWithData(FilmDeleteResponse{FilmID: filmID})
}
