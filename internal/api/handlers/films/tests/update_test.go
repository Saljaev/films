package tests

import (
	"testing"
	"tiny/internal/api/handlers/films"
)

func TestFilmsUpdateRequest_IsValid(t *testing.T) {
	type fields struct {
		Name        string
		Description string
		Rating      float64
		ReleaseDate string
		Actors      []*filmshandler.Actors
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Success",
			fields: fields{
				Name:        "test",
				Description: "about film",
				Rating:      5,
				ReleaseDate: "2020-10-10",
			},
			want: true,
		},
		{
			name: "InvalidRating",
			fields: fields{
				Name:        "test",
				Description: "about film",
				Rating:      -5,
				ReleaseDate: "2020-10-10",
			},
			want: false,
		},
		{
			name: "InvalidDate",
			fields: fields{
				Name:        "test",
				Description: "about film",
				Rating:      5,
				ReleaseDate: "2026-10-10",
			},
			want: false,
		},
		{
			name: "EmptyBody",
			fields: fields{
				Name:        "",
				Description: "",
				Rating:      0,
				ReleaseDate: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &filmshandler.FilmsUpdateRequest{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Rating:      tt.fields.Rating,
				ReleaseDate: tt.fields.ReleaseDate,
				Actors:      tt.fields.Actors,
			}
			if got := req.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
