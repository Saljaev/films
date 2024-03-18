package filmshandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tiny/internal/api/utilapi"
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
	return req.Name != "" || req.ActorName != ""
}

func GenerateResponse(films []*models.Films, res map[string]*FilmsWithActorsResponse) {
	for i := range films {
		film, ok := res[films[i].Name]

		var actor Actors

		if films[i].Actors[0].Id != 0 {
			actor = Actors{
				FirstName:   films[i].Actors[0].FirstName,
				LastName:    films[i].Actors[0].LastName,
				Gender:      films[i].Actors[0].Gender,
				DateOfBirth: films[i].Actors[0].DateOfBirth.Format(time.DateOnly),
			}
		}

		if ok {
			film.Actors[films[i].Actors[0].Id] = &actor
		} else {
			actors := make(map[int]*Actors)
			if films[i].Actors[0].Id != 0 {
				actors[films[i].Actors[0].Id] = &actor
			}

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
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	ctx.Info("request decoded", "request", req)

	films := make(map[string]*FilmsWithActorsResponse)

	if req.Name != "" {
		filmsByName, err := h.films.SearchByFilmName(context.Background(), req.Name)

		if err != nil {
			ctx.Error("failed to search film by film name", err)
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}

		GenerateResponse(filmsByName, films)
	}

	if req.ActorName != "" {
		for s := 0; s < 2; s++ {
			name := strings.Split(req.ActorName, " ")
			if len(name) == 1 {
				name = append(name, " ")
			}
			fmt.Println(name[s%2])
			filmsByActorName, err := h.films.SearchByActorName(context.Background(), name[s%2], name[s+1%2])
			if err != nil {
				ctx.Error("failed to search film by actor name", err)
				ctx.WriteFailure(http.StatusInternalServerError, "server error")
				return
			}

			GenerateResponse(filmsByActorName, films)
		}
	}

	var f []*FilmsWithActorsResponse

	for i := range films {
		f = append(f, films[i])
	}

	if len(films) == 0 {
		ctx.Info("no such films with fragment", "fragments", req)
		ctx.SuccessWithData(FilmsSearchByFragmentResponse{nil})
	} else {
		ctx.Info("successful film search", "count films", len(films))
		ctx.SuccessWithData(FilmsSearchByFragmentResponse{f})
	}
}
