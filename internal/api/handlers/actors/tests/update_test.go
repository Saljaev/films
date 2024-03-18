package tests

import (
	"testing"
	"tiny/internal/api/handlers/actors"
)

func TestActorsUpdateRequest_IsValid(t *testing.T) {
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
			name: "OnlyFistName",
			fields: fields{
				FirstName: "Johny",
			},
			want: true,
		},
		{
			name: "OnlyGender",
			fields: fields{
				Gender: "male",
			},
			want: true,
		},
		{
			name: "OnlyDate",
			fields: fields{
				DateOfBirth: "2010-10-10",
			},
			want: true,
		},
		{
			name: "InvalidGender",
			fields: fields{
				Gender: "gyro",
			},
			want: false,
		},
		{
			name: "InvalidDate",
			fields: fields{
				DateOfBirth: "2026-10-10",
			},
			want: false,
		},
		{
			name: "InvalidDate2",
			fields: fields{
				DateOfBirth: "1500-10-10",
			},
			want: false,
		},
		{
			name: "InvalidGender",
			fields: fields{
				Gender: "gyro",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &actorshandler.ActorsUpdateRequest{
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
