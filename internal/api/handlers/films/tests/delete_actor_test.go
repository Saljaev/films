package tests

import (
	"testing"
	"tiny/internal/api/handlers/films"
)

func TestFilmDeleteActorRequest_IsValid(t *testing.T) {
	type fields struct {
		FilmID  int
		ActorID int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Success",
			fields: fields{
				FilmID:  1,
				ActorID: 1,
			},
			want: true,
		},
		{
			name: "InvalidActorID",
			fields: fields{
				FilmID:  1,
				ActorID: 0,
			},
			want: false,
		},
		{
			name: "InvalidFilmID",
			fields: fields{
				FilmID:  0,
				ActorID: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &filmshandler.FilmDeleteActorRequest{
				FilmID:  tt.fields.FilmID,
				ActorID: tt.fields.ActorID,
			}
			if got := req.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
