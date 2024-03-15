package request

import (
	"time"
)

type (
	DateTime struct {
		Format string
		time.Time
	}

	UserRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	Actor struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"second_name"`
		DateOfBirth string `json:"date_of_birth"`
		Gender      string `json:"gender"`
		//Films       []string  `json:"films"`
	}

	Films struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Rating      float64  `json:"rating"`
		ReleaseDate string   `json:"release_date"`
		Actors      []*Actor `json:"actors"`
	}
)
