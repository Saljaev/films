package middleware

import (
	"context"
	"net/http"
	"tiny/internal/api/utilapi"
)

func (m *Middleware) CheckRole(ctx *utilapi.APIContext) {
	token := ctx.GetCookieString("Token")
	claims, err := m.jwt.Parse(token)
	if err != nil {
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	user, err := m.users.GetById(context.Background(), claims.Id)
	if err != nil {
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	if user.Role == "user" {
		ctx.WriteFailure(http.StatusForbidden, "forbidden")
		return
	}
}
