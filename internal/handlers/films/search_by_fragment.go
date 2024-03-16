package filmshandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"tiny/internal/api/request"
	"tiny/internal/api/response"
	"tiny/internal/logger/sl"
)

func (f *FilmsHandler) SearchByFragment(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "FilmsHandler"

		log := log.With(
			slog.String("op", op),
		)

		var req request.GetFilms

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))

			return
		}

		if req.Name == "" && req.ActorName == "" {
			log.Error("empty body")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("at least one the fields actor name of film name must be not empty"))

			return
		}

		log.Info("request decoded", slog.Any("request", req))

		films := make(map[string]*response.FilmsWithActors)

		if req.Name != "" {
			filmsByName, err := f.films.SearchByFilmName(r.Context(), req.Name)

			if err != nil {
				log.Error("failed to search film by film name", sl.Err(err))

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("internal error"))

				return
			}

			for i := range filmsByName {
				film, ok := films[filmsByName[i].Name]

				actor := response.Actors{
					FirstName:   filmsByName[i].Actors[0].FirstName,
					LastName:    filmsByName[i].Actors[0].LastName,
					Gender:      filmsByName[i].Actors[0].Gender,
					DateOfBirth: filmsByName[i].Actors[0].DateOfBirth.Format(time.DateOnly),
				}

				if ok {
					film.Actors[filmsByName[i].Actors[0].Id] = &actor
				} else {
					actors := make(map[int]*response.Actors)
					actors[filmsByName[i].Actors[0].Id] = &actor

					films[filmsByName[i].Name] = &response.FilmsWithActors{
						Name:        filmsByName[i].Name,
						Description: filmsByName[i].Description,
						Rating:      filmsByName[i].Rating,
						ReleaseDate: filmsByName[i].ReleaseDate.Format(time.DateOnly),
						Actors:      actors,
					}
				}
			}
		}

		if req.ActorName != "" {
			name := strings.Split(req.ActorName, " ")
			if len(name) > 2 {
				log.Error("to many arguments is name", slog.Any("name", name))

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("to many words in name"))

				return
			}

			for s := 0; s < 2; s++ {
				filmsByActorName, err := f.films.SearchByActorName(r.Context(), name[s%2], name[s%2])
				if err != nil {
					log.Error("failed to search film by actor name", sl.Err(err))

					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("internal error"))

					return
				}

				for i := range filmsByActorName {
					film, ok := films[filmsByActorName[i].Name]

					actor := response.Actors{
						FirstName:   filmsByActorName[i].Actors[0].FirstName,
						LastName:    filmsByActorName[i].Actors[0].LastName,
						Gender:      filmsByActorName[i].Actors[0].Gender,
						DateOfBirth: filmsByActorName[i].Actors[0].DateOfBirth.Format(time.DateOnly),
					}

					if ok {
						film.Actors[filmsByActorName[i].Actors[0].Id] = &actor
					} else {
						actors := make(map[int]*response.Actors)
						actors[filmsByActorName[i].Actors[0].Id] = &actor

						films[filmsByActorName[i].Name] = &response.FilmsWithActors{
							Name:        filmsByActorName[i].Name,
							Description: filmsByActorName[i].Description,
							Rating:      filmsByActorName[i].Rating,
							ReleaseDate: filmsByActorName[i].ReleaseDate.Format(time.DateOnly),
							Actors:      actors,
						}
					}
				}
			}
		}

		var filmArr []*response.FilmsWithActors

		for _, film := range films {
			filmArr = append(filmArr, film)
		}

		if len(filmArr) == 0 {
			log.Info("no such films with fragments", slog.Any("fragments", req))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("no films found"))
		} else {
			res, err := json.Marshal(filmArr)
			if err != nil {
				log.Error("failed to marhal films", sl.Err(err))

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal error"))
			}
			log.Info("successful film search")

			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	}
}
