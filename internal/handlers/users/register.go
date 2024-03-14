package userhandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"tiny/internal/api/request"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
)

func (h *UserHandler) Register(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "UserHandler - Register"

		log := log.With(
			slog.String("op", op),
		)

		var req request.UserRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))

			w.Write([]byte("invalid request"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		user := models.User{
			Login:    req.Login,
			Password: req.Password,
		}

		id, err := h.users.Register(r.Context(), user)
		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			w.Write([]byte("failed to register"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("user add", slog.Any("id", id))

		w.Write([]byte("user register"))
		w.WriteHeader(http.StatusOK)
	}
}
