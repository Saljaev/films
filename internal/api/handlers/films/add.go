package filmshandler

import (
	"context"
	"errors"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
	"unicode/utf8"
)

type FilmsAddRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	ReleaseDate string    `json:"release_date"`
	Actors      []*Actors `json:"actors"`
}

type FilmAddResponse struct {
	FilmId int `json:"film_id"`
}

func (a *Actors) IsValid() bool {
	if a.FirstName == "" || a.LastName == "" || a.Gender == "" || a.DateOfBirth == "" {
		return false
	}

	validGenders := map[string]struct{}{
		"male":   {},
		"female": {},
		"other":  {},
	}
	_, isValid := validGenders[a.Gender]
	_, err := time.Parse(time.DateOnly, a.DateOfBirth)

	return isValid && err == nil
}

func (req *FilmsAddRequest) IsValid() bool {
	date, err := time.Parse(time.DateOnly, req.ReleaseDate)
	return (utf8.RuneCountInString(req.Name) >= 1 && utf8.RuneCountInString(req.Name) <= 150) &&
		utf8.RuneCountInString(req.Description) <= 1000 && req.Rating > 0 && req.Rating <= 10 &&
		err == nil && date.Year() >= FilmValidDate && date.Year() <= time.Now().Year()
}

func (f *FilmsHandler) Add(ctx *utilapi.APIContext) {
	var req FilmsAddRequest
	err := ctx.Decode(&req)
	if err != nil {
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
		Name:        req.Name,
		Description: req.Description,
		Rating:      req.Rating,
		ReleaseDate: date,
		Actors:      actors,
	}

	id, err := f.films.Add(context.Background(), film)
	if err != nil {
		ctx.Error("failed to add film", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film add", "id", id)

	ctx.SuccessWithData(FilmAddResponse{FilmId: id})

}
