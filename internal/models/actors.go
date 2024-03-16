package models

import "time"

type Actor struct {
	Id          int       `json:"id,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Films       []*Films
}
