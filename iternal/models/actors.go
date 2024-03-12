package model

import "time"

type Actors struct {
	Id          int       `json: "id, omitempty"`
	FirstName   string    `json: "first_name,omitempty"`
	LastName    string    `json: "last_name,omitempty"`
	Gender      string    `json: "gender,omitempty"`
	DateOfBirth time.Time `json: "date_of_birth,omitempty"`
}
