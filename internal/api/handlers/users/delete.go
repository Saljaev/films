package userhandler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"tiny/internal/api/utilapi"
	"tiny/internal/logger/sl"
)

type UserDeleteResponse struct {
	UserID int `json:"user_id"`
}

func (h *UserHandler) Delete(ctx *utilapi.APIContext) {
	rawUserID := ctx.GetURLParam("id")

	userID, err := strconv.Atoi(rawUserID)
	if err != nil || userID <= 0 {
		ctx.Error("invalid user id from url param", sl.Err(err))
		ctx.WriteFailure(http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.users.GetById(context.Background(), userID)
	if err != nil {
		ctx.Error("failed to delete user by id", sl.Err(err), slog.Any("user_id", userID))
		ctx.WriteFailure(http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.users.Delete(context.Background(), user.Id)
	if err != nil {
		ctx.Error("failed to delete user", sl.Err(err))
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("user deleted", slog.Any("user_id", userID))
	ctx.SuccessWithData(UserDeleteResponse{UserID: userID})
}
