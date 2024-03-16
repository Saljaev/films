package filmshandler

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func (f *FilmsHandler) Delete(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "FilmsHandler - Delete"

		log := log.With(
			slog.String("op", op),
		)

		urlId := r.URL.Query().Get("id")
		if urlId == "" {
			log.Error("failed to parse id from url")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid url"))

			return
		}

		id, err := strconv.Atoi(urlId)
		if err != nil {
			log.Error("invalid id value")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid id"))

			return
		}

		log.Info("request decoded", slog.Any("id", id))

		err = f.films.Delete(r.Context(), id)
		if errors.Is(err, os.ErrNotExist) {
			log.Error("no such film with tihs id", slog.Any("id", id))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no film with given id"))

			return
		}

		if err != nil {
			log.Error("failed to delete film")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to delete film"))

			return
		}

		log.Info("successful delete")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful delete"))
	}
}
