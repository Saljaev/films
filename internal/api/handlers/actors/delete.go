package actorshandler

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"tiny/internal/logger/sl"
)

func (h *ActorsHandler) Delete(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "ActorsHandler - Delete"

		log := log.With(
			slog.String("op", op),
		)

		urlID := r.URL.Query().Get("id")
		if urlID == "" {
			log.Error("failed to parse id from url")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid url"))

			return
		}

		id, err := strconv.Atoi(urlID)
		if err != nil {
			log.Error("invalid id value")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid id"))

			return
		}

		log.Info("request decoded", slog.Any("id", id))

		err = h.actors.Delete(r.Context(), id)
		if errors.Is(err, os.ErrNotExist) {
			log.Error("no such actor with this id", slog.Any("id", id))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no actor with given id"))

			return
		}

		if err != nil {
			log.Error("failed to delete actor", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))

			return
		}

		log.Info("delete successful")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delete successful"))

	}
}
