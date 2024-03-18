package userhandler

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"tiny/internal/api/utilapi"
	"unicode/utf8"
)

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	UserID int `json:"user_id"`
}

func (req *UserLoginRequest) IsValid() bool {
	return utf8.RuneCountInString(req.Login) >= 0 && utf8.RuneCountInString(req.Password) >= 0
}

func (h *UserHandler) Login(ctx *utilapi.APIContext) {
	var req UserLoginRequest
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	user, err := h.users.GetByLogin(context.Background(), req.Login)
	if err != nil {
		ctx.Error("failed to get user by login", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.Error("failed to compare password", err)

		ctx.WriteFailure(http.StatusBadRequest, "invalid login or password")
		return
	}

	accessToken, err := h.jwt.Generate(*user, h.tokenTTL)
	if err != nil {
		ctx.Error("failed to generate access token", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	refreshToken, err := h.jwt.Generate(*user, h.sessionTTL)
	if err != nil {
		ctx.Error("failed to generate refresh token", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	sessionCheck, err := h.sessions.GetByUserId(context.Background(), user.Id)

	if errors.Is(err, os.ErrNotExist) {
		_, err = h.sessions.Add(context.Background(), refreshToken, user.Id, h.sessionTTL)
		if err != nil {
			ctx.Error("failed to create session", err)

			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}
	} else if err != nil {
		ctx.Error("failed to get session", err)

		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	} else {
		err = h.sessions.Update(context.Background(), refreshToken, sessionCheck.Id, h.sessionTTL)
		if err != nil {
			ctx.Error("failed to update session", err)

			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}
	}

	ctx.SetTokensCookie(accessToken, refreshToken, h.tokenTTL, h.sessionTTL)

	ctx.Info("login successful", "id", user.Id)

	ctx.SuccessWithData(UserLoginResponse{UserID: user.Id})

}
