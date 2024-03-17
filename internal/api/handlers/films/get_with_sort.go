package filmshandler

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
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
	ctx.Decode(&req)

	ctx.Info("request decoded", slog.Any("request", req))

	films, err := h.films.RateByField(context.Background(), req.Field, req.Increasing)
	if err != nil {
		ctx.Error("faield to get sorted films", sl.Err(err))
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

	ctx.Info("sorted films get")

	ctx.SuccessWithData(FilmsGetSortedResponse{responseFilm})
}

//func (f *FilmsHandler) GetWithSort(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "FilmsHandler - GetWithSort"
//
//		log := log.With(
//			slog.String("op", op),
//		)
//
//		var req request.GetFilmSort
//
//		err := json.NewDecoder(r.Body).Decode(&req)
//
//		if err != nil {
//			log.Error("failed to decode body", sl.Err(err))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid request"))
//
//			return
//		}
//
//		log.Info("request decoded", slog.Any("request", req))
//
//		films, err := f.films.RateByField(r.Context(), req.Field, req.Increasing)
//		if err != nil {
//			log.Error("failed to get sorted films", sl.Err(err))
//
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte("internal error"))
//
//			return
//		}
//
//		var res []*response.Films
//		for i := range films {
//			film := response.Films{
//				Name:        films[i].Name,
//				Description: films[i].Description,
//				Rating:      films[i].Rating,
//				ReleaseDate: films[i].ReleaseDate.Format(time.DateOnly),
//			}
//
//			res = append(res, &film)
//		}
//
//		data, err := json.Marshal(res)
//		if err != nil {
//			log.Error("failed to marshal data", sl.Err(err))
//
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte("internal error"))
//
//			return
//		}
//
//		log.Info("successful get")
//		w.WriteHeader(http.StatusOK)
//		w.Write(data)
//	}
//}
