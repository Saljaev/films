package actorshandler

import (
	"log/slog"
	"net/http"
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

		//TODO: add os.ErrNoExists to deleted data

		err = h.actors.Delete(r.Context(), id)
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
