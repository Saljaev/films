package filmshandler

import (
	"context"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
)

type FilmsGetSortedRequest struct {
	Field      string `json:"field"`
	Increasing bool   `json:"increasing"`
}

type FilmsGetSortedResponse struct {
	Films []*Films `json:"films"`
}

func (req *FilmsGetSortedRequest) IsValid() bool {
	fields := map[string]struct{}{
		"name":           {},
		"rating":         {},
		"releasing_date": {},
	}
	_, ok := fields[req.Field]
	if !ok {
		req.Field = "rating"
	}

	return true
}

func (h *FilmsHandler) GetWithSort(ctx *utilapi.APIContext) {
	var req FilmsGetSortedRequest
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	ctx.Info("request decoded", "request", req)

	films, err := h.films.RateByField(context.Background(), req.Field, req.Increasing)
	if err != nil {
		ctx.Error("faield to get sorted films", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	var responseFilm []*Films

	for i := range films {
		film := Films{
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			ReleaseDate: films[i].ReleaseDate.Format(time.DateOnly),
		}

		responseFilm = append(responseFilm, &film)
	}

	ctx.Info("sorted films get", "count films", len(responseFilm))

	ctx.SuccessWithData(FilmsGetSortedResponse{responseFilm})
}
