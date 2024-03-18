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

const ActorValidDate = 1800

func ValidateDate(dateOfBirth string) bool {
	date, err := time.Parse(time.DateOnly, dateOfBirth)
	if err != nil {
		return false
	} else if date.Year() < ActorValidDate || date.Year() > time.Now().Year() {
		return false
	}

	return true

}

type Films struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"release_date"`
}
