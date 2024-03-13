package entities

type Films struct {
	Id          int
	Name        string
	Description string
	Rating      float64
	Actors      []*Actors
}
