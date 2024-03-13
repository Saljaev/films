package userhandler

import (
	"log/slog"
	"net/http"
	"tiny/internal/logger/sl"
)

// Delete user by login in URL parameters
func (h *UserHandler) Delete(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Userhandler - Delete"

		log := log.With(
			slog.String("op", op),
		)

		// TODO: change login to id

		login := r.URL.Query().Get("login")
		if login == "" {
			log.Error("failed to parse login from url")

			w.Write([]byte("invalid url"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		_, err := h.users.GetByLogin(r.Context(), login)
		if err != nil {
			log.Error("failed to delete there is no user with login:"+login, sl.Err(err))

			w.Write([]byte("invalid login"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = h.users.Delete(r.Context(), login)
		if err != nil {
			log.Error("failed to delete user", sl.Err(err))

			w.Write([]byte("internal error"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("user deleted", slog.Any("login", login))

		w.Write([]byte("delete successful"))
		w.WriteHeader(http.StatusOK)
	}
}
