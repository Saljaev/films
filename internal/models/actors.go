package models

import "time"

type Actors struct {
	Id          int
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	Films       []*Films
}
