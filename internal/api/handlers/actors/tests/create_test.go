package tests

import (
	"testing"
	"tiny/internal/api/handlers/actors"
)

func TestActorCreate_IsValid(t *testing.T) {
	type fields struct {
		FirstName   string
		LastName    string
		DateOfBirth string
		Gender      string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Success",
			fields: fields{
				FirstName:   "Johny",
				LastName:    "Depp",
				DateOfBirth: "2023-01-01",
				Gender:      "male",
			},
			want: true,
		},
		{
			name: "EmptyFirstName",
			fields: fields{
				FirstName:   "",
				LastName:    "Depp",
				DateOfBirth: "2010-10-10",
				Gender:      "male",
			},
			want: false,
		},
		{
			name: "EmptyLastName",
			fields: fields{
				FirstName:   "Johny",
				LastName:    "",
				DateOfBirth: "2010-10-10",
				Gender:      "male",
			},
			want: false,
		},
		{
			name: "InvalidDate",
			fields: fields{
				FirstName:   "Johny",
				LastName:    "Depp",
				DateOfBirth: "2026-10-10",
				Gender:      "male",
			},
			want: false,
		},
		{
			name: "InvalidDate2",
			fields: fields{
				FirstName:   "Johny",
				LastName:    "Depp",
				DateOfBirth: "1500-10-10",
				Gender:      "male",
			},
			want: false,
		},
		{
			name: "InvalidGender",
			fields: fields{
				FirstName:   "Johny",
				LastName:    "Depp",
				DateOfBirth: "2010-10-10",
				Gender:      "gyro",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &actorshandler.ActorCreate{
				FirstName:   tt.fields.FirstName,
				LastName:    tt.fields.LastName,
				DateOfBirth: tt.fields.DateOfBirth,
				Gender:      tt.fields.Gender,
			}
			if got := req.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
