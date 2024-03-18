package filmshandler

import "tiny/internal/usecase"

type FilmsHandler struct {
	films usecase.Films
}

func NewFilmsHandler(f usecase.Films) *FilmsHandler {
	return &FilmsHandler{f}
}

const FilmValidDate = 1500

type Films struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"release_date"`
}

type Actors struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}
