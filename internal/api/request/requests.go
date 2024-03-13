package request

import "time"

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Actor struct {
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"second_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}
