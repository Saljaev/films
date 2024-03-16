package entities

import "time"

type Films struct {
	Id          int
	Name        string
	Description string
	Rating      float64
	ReleaseDate time.Time
	Actors      []*Actors
}
