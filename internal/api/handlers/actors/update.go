package actorshandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"tiny/internal/api/request"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
)

func (h *ActorsHandler) Update(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "ActorsHandler - Update"

		log := log.With(
			slog.String("op", op),
		)

		var req request.Actor

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			log.Error("failed to decode body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))

			return
		}

		if req.FirstName == "" && req.LastName == "" && req.Gender == "" && req.DateOfBirth == "" {
			log.Error("empty body")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("enter at least one field to update"))

			return
		}

		if !validateGender(req.Gender) {
			log.Error("invalid actor's gender", slog.Any("gender", req.Gender))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("enter actor's gender from the available one: [male/female/other]"))

			return
		}

		date, err := time.Parse(time.DateOnly, req.DateOfBirth)
		if err != nil {
			log.Error("invalid date format", slog.Any("date", req.DateOfBirth))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("enter actor's date of birth in format YYYY-MM-DD"))

			return
		}

		log.Info("request decoded", slog.Any("request", req))

		actor := models.Actor{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      req.Gender,
			DateOfBirth: date,
		}

		err = h.actors.Update(r.Context(), actor)
		if err != nil {
			log.Error("failed to update actor", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to update actor"))

			return
		}

		log.Info("successful update")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user updated"))
	}
}
