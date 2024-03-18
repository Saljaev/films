package userhandler

import (
	"context"
	"net/http"
	"tiny/internal/api/utilapi"
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
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	user := models.User{
		Login:    req.Login,
		Password: req.Password,
	}

	id, err := h.users.Register(context.Background(), user)
	if err != nil {
		ctx.Error("failed to create user", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("user add", "id", id)

	ctx.SuccessWithData(UserRegisterResponse{UserID: id})
}
