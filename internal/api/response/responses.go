package response

import "time"

type IdResponse struct {
	Id int `json:"id"`
}

type Actors struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
}
