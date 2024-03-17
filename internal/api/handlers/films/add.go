package filmshandler

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
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
		err == nil && date.Year() >= 1700
}

func (f *FilmsHandler) Add(ctx *utilapi.APIContext) {
	var req FilmsAddRequest
	ctx.Decode(&req)

	var actors []*models.Actors

	for i := range req.Actors {
		if !req.Actors[i].IsValid() {
			ctx.Error("failed to decode body")
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

	ctx.Info("request decoded", slog.Any("request", req))

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
		ctx.Error("failed to add film", sl.Err(err))

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film add", slog.Any("id", id))

	ctx.SuccessWithData(FilmAddResponse{FilmId: id})

}

//func (f *FilmsHandler) Add(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "FilmsHandler - Add"
//
//		log := log.With(
//			slog.String("op", op),
//		)
//
//		var req request.Films
//
//		err := json.NewDecoder(r.Body).Decode(&req)
//
//		if err != nil {
//			log.Error("failed to decode body", sl.Err(err))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid request"))
//
//			return
//		}
//
//		if len(req.Name) < 1 || len(req.Name) >= 150 {
//			log.Error("invalid film's name", slog.Any("name", req.Name))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("film's name must be at least 1 characters or no more than 150 characters"))
//
//			return
//		}
//
//		if len(req.Description) > 1000 {
//			log.Error("invalid film's description", slog.Any("description", req.Description))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("film's description must be no more than 1000 characters"))
//
//			return
//		}
//
//		if req.Rating < 0 || req.Rating > 10 {
//			log.Error("invalid film's rating", slog.Any("rating", req.Rating))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("film's rating must be at least 0 or no more than 10"))
//
//			return
//		}
//
//		date, err := time.Parse(time.DateOnly, req.ReleaseDate)
//
//		if err != nil {
//			log.Error("invalid film's date format", slog.Any("date", req.ReleaseDate))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("enter release date in format YYYY-MM-DD"))
//
//			return
//		}
//
//		var actors []*models.Actors
//
//		for i := range req.Actors {
//			time, err := time.Parse(time.DateOnly, req.Actors[i].DateOfBirth)
//			if err != nil {
//				log.Error("invalid film actor's date format", slog.Any("date", req.Actors[i].DateOfBirth))
//
//				w.WriteHeader(http.StatusBadRequest)
//				w.Write([]byte("enter actor's date of birth in format YYYY-MM-DD"))
//
//				return
//			}
//
//			isGenderValid := ValidateGender(req.Actors[i].Gender)
//			if !isGenderValid {
//				log.Error("invalid film actor's gender", slog.Any("gender", req.Actors[i].Gender))
//
//				w.WriteHeader(http.StatusBadRequest)
//				w.Write([]byte("enter actor's gender from the available one: [male/female/other]"))
//
//				return
//			}
//
//			actor := models.Actors{
//				FirstName:   req.Actors[i].FirstName,
//				LastName:    req.Actors[i].LastName,
//				Gender:      req.Actors[i].Gender,
//				DateOfBirth: time,
//			}
//
//			actors = append(actors, &actor)
//		}
//
//		log.Info("request decoded", slog.Any("request", req))
//
//		film := models.Films{
//			Name:        req.Name,
//			Description: req.Description,
//			Rating:      req.Rating,
//			ReleaseDate: date,
//			Actors:      actors,
//		}
//
//		id, err := f.films.Add(r.Context(), film)
//		if err != nil {
//			log.Error("failed to add film", sl.Err(err))
//
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte("failed to add film"))
//
//			return
//		}
//
//		log.Info("film add", slog.Any("id", id))
//
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("film add"))
//	}
//}
//
//func ValidateGender(gender string) bool {
//	validGenders := map[string]bool{
//		"male":   true,
//		"female": true,
//		"other":  true,
//	}
//	_, isValid := validGenders[gender]
//	return isValid
//}