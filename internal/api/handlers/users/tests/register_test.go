package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiny/internal/api/handlers/users"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
	"tiny/internal/usecase/mocks"
	"tiny/pkg"
)

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name       string
		login      string
		password   string
		respId     int
		statusCode int
	}{
		{
			name:       "Success",
			login:      "test",
			password:   "qwerty123",
			respId:     1,
			statusCode: http.StatusOK,
		},
		{
			name:       "ShortPassword",
			login:      "test",
			password:   "qwe",
			respId:     0,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "EmptyLogin",
			login:      "",
			password:   "qwerty123",
			respId:     0,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := mocks.NewUsers(t)

			h := &userhandler.UserHandler{
				users: users,
			}

			if tt.respId != 0 {
				users.On("Register", context.Background(), models.User{Login: tt.login, Password: tt.password}).Return(1, nil).Once()
			}

			log := pkg.SetupPrettySlog()

			router := utilapi.NewRouter(log)
			router.Handle("/", h.Register)

			input := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, tt.login, tt.password)
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tt.statusCode)
		})
	}
}

func TestUserRegisterRequest_IsValid(t *testing.T) {
	type fields struct {
		Login    string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Success",
			fields: fields{
				Login:    "test",
				Password: "qwerty123",
			},
			want: true,
		},
		{
			name: "LoginEmpty",
			fields: fields{
				Login:    "",
				Password: "qwerty123",
			},
			want: false,
		},
		{
			name: "PasswordEmpty",
			fields: fields{
				Login:    "test",
				Password: "",
			},
			want: false,
		},
		{
			name: "PasswordToShort",
			fields: fields{
				Login:    "test",
				Password: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &userhandler.UserRegisterRequest{
				Login:    tt.fields.Login,
				Password: tt.fields.Password,
			}
			if got := req.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
