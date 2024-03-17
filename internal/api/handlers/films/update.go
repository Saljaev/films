package filmshandler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
	"unicode/utf8"
)

type FilmsUpdateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Rating      float64  `json:"rating"`
	ReleaseDate string   `json:"release_date"`
	Actors      []*Actor `json:"actors"`
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
	if req.ReleaseDate != "" {
		date, err := time.Parse(time.DateOnly, req.ReleaseDate)
		if err != nil || date.Year() <= 1700 {
			return false
		}
	}

	return true
}

func (h *FilmsHandler) Upadate(ctx *utilapi.APIContext) {
	var req FilmsUpdateRequest
	ctx.Decode(&req)

	rawFilmID := ctx.GetURLParam("id")

	filmID, err := strconv.Atoi(rawFilmID)
	if err != nil || filmID <= 0 {
		ctx.Error("invalid film id from url param", sl.Err(err))
		ctx.WriteFailure(http.StatusBadRequest, "invalid film id")
		return
	}

	var actors []*models.Actor

	for i := range req.Actors {
		if !req.Actors[i].IsValid() {
			ctx.Error("failed to decode body")
			ctx.WriteFailure(http.StatusBadRequest, "server error")
			return
		}

		date, _ := time.Parse(time.DateOnly, req.Actors[i].DateOfBirth)

		actor := models.Actor{
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

	err = h.films.Update(context.Background(), film)
	if err != nil {
		ctx.Error("failed to update film", sl.Err(err))
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("film update", slog.Any("id", filmID))
	ctx.SuccessWithData(FilmAddResponse{FilmId: filmID})
}

//func (f *FilmsHandler) Update(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "FilmsHandler - Update"
//
//		log := log.With(
//			slog.String("op", op),
//		)
//
//		var req request.Films
//
//		err := json.NewDecoder(r.Body).Decode(&req)
//		urlId := r.URL.Query().Get("id")
//		if urlId == "" {
//			log.Error("failed to parse id from url")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid url"))
//
//			return
//		}
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
//		id, err := strconv.Atoi(urlId)
//		if err != nil {
//			log.Error("invalid id value")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid id"))
//
//			return
//		}
//
//		if req.Name == "" && req.Rating == 0 && req.Description == "" && req.ReleaseDate == "" && req.Actors == nil {
//			log.Error("body empty")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("body empty"))
//
//			return
//		}
//
//		if req.Name != "" && len(req.Name) < 1 || len(req.Name) >= 150 {
//			log.Error("invalid film's name", slog.Any("name", req.Name))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("film's name must be at least 1 characters or no more than 150 characters"))
//
//			return
//		}
//
//		if req.Description != "" && len(req.Description) > 1000 {
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
//		var date time.Time
//
//		if req.ReleaseDate != "" {
//			date, err = time.Parse(time.DateOnly, req.ReleaseDate)
//			if err != nil {
//				log.Error("invalid film's date format", slog.Any("date", req.ReleaseDate))
//
//				w.WriteHeader(http.StatusBadRequest)
//				w.Write([]byte("enter release date in format YYYY-MM-DD"))
//
//				return
//			}
//		}
//
//		//TODO: delete actor from film
//
//		var actors []*models.Actor
//
//		if req.Actors != nil {
//			for i := range req.Actors {
//				time, err := time.Parse(time.DateOnly, req.Actors[i].DateOfBirth)
//				if err != nil {
//					log.Error("invalid film actor's date format", slog.Any("date", req.Actors[i].DateOfBirth))
//
//					w.WriteHeader(http.StatusBadRequest)
//					w.Write([]byte("enter actor's date of birth in format YYYY-MM-DD"))
//
//					return
//				}
//
//				isGenderValid := ValidateGender(req.Actors[i].Gender)
//				if !isGenderValid {
//					log.Error("invalid film actor's gender", slog.Any("gender", req.Actors[i].Gender))
//
//					w.WriteHeader(http.StatusBadRequest)
//					w.Write([]byte("enter actor's gender from the available one: [male/female/other]"))
//
//					return
//				}
//
//				actor := models.Actor{
//					FirstName:   req.Actors[i].FirstName,
//					LastName:    req.Actors[i].LastName,
//					Gender:      req.Actors[i].Gender,
//					DateOfBirth: time,
//				}
//
//				actors = append(actors, &actor)
//			}
//
//		}
//
//		log.Info("request decoded", slog.Any("request", req))
//
//		film := models.Films{
//			Id:          id,
//			Name:        req.Name,
//			Description: req.Description,
//			Rating:      req.Rating,
//			ReleaseDate: date,
//			Actors:      actors,
//		}
//
//		err = f.films.Update(r.Context(), film)
//		if err != nil {
//			log.Error("failed to update film", sl.Err(err))
//
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte("failed to update film"))
//
//			return
//		}
//
//		log.Info("film update successful")
//
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("update successful"))
//	}
//}
