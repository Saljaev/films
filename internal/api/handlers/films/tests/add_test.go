package tests

import (
	"testing"
	"tiny/internal/api/handlers/films"
)

func TestFilmsAddRequest_IsValid(t *testing.T) {
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
				Description: "Films about",
				Rating:      1,
				ReleaseDate: "2023-10-10",
			},
			want: true,
		},
		{
			name: "NegativeRating",
			fields: fields{
				Name:        "test",
				Description: "Films about",
				Rating:      -1,
				ReleaseDate: "2000-10-10",
			},
			want: false,
		},
		{
			name: "EmptyName",
			fields: fields{
				Name:        "",
				Description: "Films about",
				Rating:      1,
				ReleaseDate: "2000-10-10",
			},
			want: false,
		},
		{
			name: "InvalidDate",
			fields: fields{
				Name:        "test",
				Description: "Films about",
				Rating:      1,
				ReleaseDate: "2026-10-10",
			},
			want: false,
		},
		{
			name: "InvalidDat2",
			fields: fields{
				Name:        "test",
				Description: "Films about",
				Rating:      1,
				ReleaseDate: "1400-10-10",
			},
			want: false,
		},
		{
			name: "SuccessWithActor",
			fields: fields{
				Name:        "test",
				Description: "Films about",
				Rating:      1,
				ReleaseDate: "2020-10-10",
				Actors: []*filmshandler.Actors{
					{
						FirstName:   "Johny",
						LastName:    "Depp",
						Gender:      "male",
						DateOfBirth: "2010-10-10",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &filmshandler.FilmsAddRequest{
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
