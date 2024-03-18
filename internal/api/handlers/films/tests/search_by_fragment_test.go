package tests

import (
	"testing"
	"tiny/internal/api/handlers/films"
)

func TestFilmsSearchByFragmentRequest_IsValid(t *testing.T) {
	type fields struct {
		Name      string
		ActorName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Success",
			fields: fields{
				Name:      "test",
				ActorName: "Johny Depp",
			},
			want: true,
		},
		{
			name: "EmptyRequest",
			fields: fields{
				Name:      "",
				ActorName: "",
			},
			want: false,
		},
		{
			name: "OnlyActorName",
			fields: fields{
				Name:      "",
				ActorName: "Johny Depp",
			},
			want: true,
		},
		{
			name: "OnlyFilmName",
			fields: fields{
				Name:      "test",
				ActorName: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &filmshandler.FilmsSearchByFragmentRequest{
				Name:      tt.fields.Name,
				ActorName: tt.fields.ActorName,
			}
			if got := req.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
