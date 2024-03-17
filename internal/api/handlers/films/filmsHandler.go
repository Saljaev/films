package filmshandler

import "tiny/internal/usecase"

type FilmsHandler struct {
	films usecase.Films
}

func NewFilmsHandler(f usecase.Films) *FilmsHandler {
	return &FilmsHandler{f}
}
