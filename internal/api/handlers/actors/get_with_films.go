package actorshandler

import (
	"context"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
)

type (
	ActorsWithFilmsResponse struct {
		Actors []*ActorsAndFilmsResponse
	}

	ActorsAndFilmsResponse struct {
		FirstName   string   `json:"first_name"`
		LastName    string   `json:"last_name"`
		Gender      string   `json:"gender"`
		DateOfBirth string   `json:"date_of_birth"`
		Films       []*Films `json:"films"`
	}
)

func GenerateResponse(actors []*models.Actors) []*ActorsAndFilmsResponse {
	var res []*ActorsAndFilmsResponse

	for i := range actors {

		film := Films{
			Name:        actors[i].Films[0].Name,
			Description: actors[i].Films[0].Description,
			Rating:      actors[i].Films[0].Rating,
			ReleaseDate: actors[i].Films[0].ReleaseDate.Format(time.DateOnly),
		}

		actor := ActorsAndFilmsResponse{
			FirstName:   actors[i].FirstName,
			LastName:    actors[i].LastName,
			DateOfBirth: actors[i].DateOfBirth.Format(time.DateOnly),
			Gender:      actors[i].Gender,
			Films:       []*Films{&film},
		}

		if film.Name == "" {
			actor.Films = actor.Films[:len(actor.Films)-1]
		}

		if len(res) > 0 {
			if res[len(res)-1].FirstName == actor.FirstName && res[len(res)-1].LastName == actor.LastName &&
				res[len(res)-1].Gender == actor.Gender && res[len(res)-1].DateOfBirth == actor.DateOfBirth {
				res[len(res)-1].Films = append(res[len(res)-1].Films, &film)
			} else {
				res = append(res, &actor)
			}
		} else {

			res = append(res, &actor)

		}

	}

	return res
}

func (h *ActorsHandler) GetWithFilms(ctx *utilapi.APIContext) {
	actorsWithFilm, err := h.actors.GetAll(context.Background())

	if err != nil {
		ctx.Error("failed to get actors with films", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	actors := GenerateResponse(actorsWithFilm)

	ctx.Info("successful get actors", "count actors", len(actors))
	ctx.SuccessWithData(ActorsWithFilmsResponse{Actors: actors})
}
