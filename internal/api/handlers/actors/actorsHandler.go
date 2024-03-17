package actorshandler

import (
	"time"
	"tiny/internal/usecase"
)

type ActorsHandler struct {
	actors usecase.Actors
}

func NewActorsHandler(a usecase.Actors) *ActorsHandler {
	return &ActorsHandler{a}
}

func ValidateGender(gender string) bool {
	validGenders := map[string]struct{}{
		"male":   {},
		"female": {},
		"other":  {},
	}
	_, isValid := validGenders[gender]
	return isValid
}

func ValidateDate(dateOfBirth string) bool {
	date, err := time.Parse(time.DateOnly, dateOfBirth)
	return err != nil || date.Year() > 1800 || date.Year() < time.Now().Year()
}

type Films struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"release_date"`
}
