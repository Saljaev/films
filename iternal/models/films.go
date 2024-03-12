package model

type Films struct {
	Id          int      `json: "id,omitempty"`
	Name        string   `json: "name,omitempty"`
	Description string   `json: "description,omitempty"`
	Rating      float64  `json: "rating,omitempty"`
	Actors      []Actors `json: "actors,omitempty`
}
