package response

type IdResponse struct {
	Id int `json:"id"`
}

type (
	Actors struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Gender      string `json:"gender"`
		DateOfBirth string `json:"date_of_birth"`
	}
	FilmsWithActors struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Rating      float64         `json:"rating"`
		ReleaseDate string          `json:"release_date"`
		Actors      map[int]*Actors `json:"actors"`
	}
	Films struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Rating      float64 `json:"rating"`
		ReleaseDate string  `json:"release_date"`
	}
)
