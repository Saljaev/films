package filmshandler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
)

type (
	FilmsSearchByFragmentRequest struct {
		Name      string `json:"name"`
		ActorName string `json:"actor_name"`
	}

	FilmsSearchByFragmentResponse struct {
		Films []*FilmsWithActorsResponse `json:"films"`
	}

	FilmsWithActorsResponse struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Rating      float64         `json:"rating"`
		ReleaseDate string          `json:"release_date"`
		Actors      map[int]*Actors `json:"actors"`
	}
)

func (req *FilmsSearchByFragmentRequest) IsValid() bool {
	name := strings.Split(req.ActorName, " ")
	return (req.ActorName != "" && len(name) <= 2) || req.Name != ""
}

func GenerateResponse(films []*models.Films, res map[string]*FilmsWithActorsResponse) {
	for i := range films {
		film, ok := res[films[i].Name]

		actor := Actors{
			FirstName:   films[i].Actors[0].FirstName,
			LastName:    films[i].Actors[0].LastName,
			Gender:      films[i].Actors[0].Gender,
			DateOfBirth: films[i].Actors[0].DateOfBirth.Format(time.DateOnly),
		}

		if ok {
			film.Actors[films[i].Actors[0].Id] = &actor
		} else {
			actors := make(map[int]*Actors)
			actors[films[i].Actors[0].Id] = &actor

			res[films[i].Name] = &FilmsWithActorsResponse{
				Name:        films[i].Name,
				Description: films[i].Description,
				Rating:      films[i].Rating,
				ReleaseDate: films[i].ReleaseDate.Format(time.DateOnly),
				Actors:      actors,
			}
		}
	}
}

func (h *FilmsHandler) SearchByFragment(ctx *utilapi.APIContext) {
	var req FilmsSearchByFragmentRequest
	ctx.Decode(&req)

	ctx.Info("request decoded", slog.Any("request", req))

	films := make(map[string]*FilmsWithActorsResponse)

	if req.Name != "" {
		filmsByName, err := h.films.SearchByFilmName(context.Background(), req.Name)

		if err != nil {
			ctx.Error("failed to search film by film name", sl.Err(err))
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}

		GenerateResponse(filmsByName, films)
	}

	if req.ActorName != "" {
		for s := 0; s < 2; s++ {
			name := strings.Split(req.ActorName, " ")
			filmsByActorName, err := h.films.SearchByActorName(context.Background(), name[s%2], name[s%2])
			if err != nil {
				ctx.Error("failed to search film by actor name", sl.Err(err))
				ctx.WriteFailure(http.StatusInternalServerError, "server error")
				return
			}

			GenerateResponse(filmsByActorName, films)
		}
	}

	var res *FilmsSearchByFragmentResponse

	for _, v := range films {
		res.Films = append(res.Films, v)
	}

	if len(films) == 0 {
		ctx.Info("no such films with fragment", slog.Any("fragments", req))
		ctx.SuccessWithData(FilmsSearchByFragmentResponse{nil})
	} else {
		ctx.Info("successful film search")
		ctx.SuccessWithData(FilmsSearchByFragmentResponse{res.Films})
	}
}

//func (f *FilmsHandler) SearchByFragment(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "FilmsHandler"
//
//		log := log.With(
//			slog.String("op", op),
//		)
//
//		var req request.GetFilms
//
//		err := json.NewDecoder(r.Body).Decode(&req)
//		if err != nil {
//			log.Error("failed to decode body", sl.Err(err))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid request"))
//
//			return
//		}
//
//		if req.Name == "" && req.ActorName == "" {
//			log.Error("empty body")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("at least one the fields actor name of film name must be not empty"))
//
//			return
//		}
//
//		log.Info("request decoded", slog.Any("request", req))
//
//		films := make(map[string]*response.FilmsWithActors)
//
//		if req.Name != "" {
//			filmsByName, err := f.films.SearchByFilmName(r.Context(), req.Name)
//
//			if err != nil {
//				log.Error("failed to search film by film name", sl.Err(err))
//
//				w.WriteHeader(http.StatusInternalServerError)
//				w.Write([]byte("internal error"))
//
//				return
//			}
//
//			for i := range filmsByName {
//				film, ok := films[filmsByName[i].Name]
//
//				actor := response.Actors{
//					FirstName:   filmsByName[i].Actors[0].FirstName,
//					LastName:    filmsByName[i].Actors[0].LastName,
//					Gender:      filmsByName[i].Actors[0].Gender,
//					DateOfBirth: filmsByName[i].Actors[0].DateOfBirth.Format(time.DateOnly),
//				}
//
//				if ok {
//					film.Actors[filmsByName[i].Actors[0].Id] = &actor
//				} else {
//					actors := make(map[int]*response.Actors)
//					actors[filmsByName[i].Actors[0].Id] = &actor
//
//					films[filmsByName[i].Name] = &response.FilmsWithActors{
//						Name:        filmsByName[i].Name,
//						Description: filmsByName[i].Description,
//						Rating:      filmsByName[i].Rating,
//						ReleaseDate: filmsByName[i].ReleaseDate.Format(time.DateOnly),
//						Actors:      actors,
//					}
//				}
//			}
//		}
//
//		if req.ActorName != "" {
//			name := strings.Split(req.ActorName, " ")
//			if len(name) > 2 {
//				log.Error("to many arguments is name", slog.Any("name", name))
//
//				w.WriteHeader(http.StatusBadRequest)
//				w.Write([]byte("to many words in name"))
//
//				return
//			}
//
//			for s := 0; s < 2; s++ {
//				filmsByActorName, err := f.films.SearchByActorName(r.Context(), name[s%2], name[s%2])
//				if err != nil {
//					log.Error("failed to search film by actor name", sl.Err(err))
//
//					w.WriteHeader(http.StatusBadRequest)
//					w.Write([]byte("internal error"))
//
//					return
//				}
//
//				for i := range filmsByActorName {
//					film, ok := films[filmsByActorName[i].Name]
//
//					actor := response.Actors{
//						FirstName:   filmsByActorName[i].Actors[0].FirstName,
//						LastName:    filmsByActorName[i].Actors[0].LastName,
//						Gender:      filmsByActorName[i].Actors[0].Gender,
//						DateOfBirth: filmsByActorName[i].Actors[0].DateOfBirth.Format(time.DateOnly),
//					}
//
//					if ok {
//						film.Actors[filmsByActorName[i].Actors[0].Id] = &actor
//					} else {
//						actors := make(map[int]*response.Actors)
//						actors[filmsByActorName[i].Actors[0].Id] = &actor
//
//						films[filmsByActorName[i].Name] = &response.FilmsWithActors{
//							Name:        filmsByActorName[i].Name,
//							Description: filmsByActorName[i].Description,
//							Rating:      filmsByActorName[i].Rating,
//							ReleaseDate: filmsByActorName[i].ReleaseDate.Format(time.DateOnly),
//							Actors:      actors,
//						}
//					}
//				}
//			}
//		}
//
//		var filmArr []*response.FilmsWithActors
//
//		for _, film := range films {
//			filmArr = append(filmArr, film)
//		}
//
//		if len(filmArr) == 0 {
//			log.Info("no such films with fragments", slog.Any("fragments", req))
//
//			w.WriteHeader(http.StatusOK)
//			w.Write([]byte("no films found"))
//		} else {
//			res, err := json.Marshal(filmArr)
//			if err != nil {
//				log.Error("failed to marhal films", sl.Err(err))
//
//				w.WriteHeader(http.StatusInternalServerError)
//				w.Write([]byte("internal error"))
//			}
//			log.Info("successful film search")
//
//			w.WriteHeader(http.StatusOK)
//			w.Write(res)
//		}
//	}
//}
