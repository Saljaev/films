package filmshandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
	"unicode/utf8"
)

type FilmsUpdateRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	ReleaseDate string    `json:"release_date"`
	Actors      []*Actors `json:"actors"`
}

type FilmsUpdateResponse struct {
	FilmID int `json:"film_id"`
}

func (req *FilmsUpdateRequest) IsValid() bool {
	if req.Name == "" && req.Rating == 0 && req.Description == "" &&
		req.ReleaseDate == "" && req.Actors == nil {
		return false
	}
	if req.Name != "" && (utf8.RuneCountInString(req.Name) < 1 || utf8.RuneCountInString(req.Name) > 150) {
		return false
	}
	if req.Description != "" && utf8.RuneCountInString(req.Description) > 1000 {
		return false
	}
	if req.Rating != 0 && req.Rating < 0 || req.Rating > 10 {
		return false
	}
	// TODO: add validate time
	if req.ReleaseDate != "" {
		date, err := time.Parse(time.DateOnly, req.ReleaseDate)
		if err != nil || date.Year() <= FilmValidDate {
			return false
		}
	}

	return true
}

func (h *FilmsHandler) Update(ctx *utilapi.APIContext) {
	var req FilmsUpdateRequest
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	rawFilmID := ctx.GetURLParam("id")

	filmID, err := strconv.Atoi(rawFilmID)
	if err != nil || filmID <= 0 {
		ctx.Error("invalid film id from url param", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid film id")
		return
	}

	var actors []*models.Actors

	for i := range req.Actors {
		if !req.Actors[i].IsValid() {
			ctx.Error("failed to decode body", errors.New("invalid request"))
			ctx.WriteFailure(http.StatusBadRequest, "server error")
			return
		}

		date, _ := time.Parse(time.DateOnly, req.Actors[i].DateOfBirth)

		actor := models.Actors{
			FirstName:   req.Actors[i].FirstName,
			LastName:    req.Actors[i].LastName,
			Gender:      req.Actors[i].Gender,
			DateOfBirth: date,
		}

		actors = append(actors, &actor)
	}

	ctx.Info("request decoded", "request", req)

	date, _ := time.Parse(time.DateOnly, req.ReleaseDate)

	film := models.Films{
		Id:          filmID,
		Name:        req.Name,
		Description: req.Description,
		Rating:      req.Rating,
		ReleaseDate: date,
		Actors:      actors,
	}

	err = h.films.Update(context.Background(), film)
	if err != nil {
		ctx.Error("failed to update film", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film update", "id", filmID)
	ctx.SuccessWithData(FilmAddResponse{FilmId: filmID})
}
