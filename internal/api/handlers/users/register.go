package userhandler

import (
	"context"
	"log/slog"
	"net/http"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
	"unicode/utf8"
)

type UserRegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (req *UserRegisterRequest) IsValid() bool {
	return utf8.RuneCountInString(req.Login) >= 2 && utf8.RuneCountInString(req.Password) >= 6
}

type UserRegisterResponse struct {
	UserID int `json:"user_id"`
}

func (h *UserHandler) Register(ctx *utilapi.APIContext) {
	var req UserRegisterRequest
	ctx.Decode(&req)

	user := models.User{
		Login:    req.Login,
		Password: req.Password,
	}

	id, err := h.users.Register(context.Background(), user)
	if err != nil {
		ctx.Error("failed to create user", sl.Err(err))

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("user add", slog.Any("id", id))

	ctx.SuccessWithData(UserRegisterResponse{UserID: id})
}

//func (h *UserHandler) Register(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "UserHandler - Register"
//
//		log := log.With(
//			slog.String("op", op),
//		)
//
//		var req request.UserRequest
//
//		err := json.NewDecoder(r.Body).Decode(&req)
//		if err != nil {
//			log.Error("failed to decode body", sl.Err(err))
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("invalid request"))
//
//			return
//		}
//
//		if req.Login == "" {
//			log.Error("try to login without username")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("field login is required"))
//
//			return
//		}
//
//		if req.Password == "" {
//			log.Error("try to login without password")
//
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte("field password is required"))
//
//			return
//		}
//
//		user := models.User{
//			Login:    req.Login,
//			Password: req.Password,
//		}
//
//		id, err := h.users.Register(r.Context(), user)
//		if err != nil {
//			log.Error("failed to create user", sl.Err(err))
//
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte("failed to register"))
//
//			return
//		}
//
//		log.Info("user add", slog.Any("id", id))
//
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("user register"))
//	}
//}
