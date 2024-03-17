package middleware

import (
	"context"
	"net/http"
	"tiny/internal/api/utilapi"
)

func (m *Middleware) Validate(ctx *utilapi.APIContext) {
	token := ctx.GetCookieString("Token")

	if token != "" {
		claims, err := m.jwt.Parse(token)
		if err != nil {
			ctx.WriteFailure(http.StatusBadRequest, "invalid access token")
			return
		}

		_, err = m.users.GetById(context.Background(), claims.Id)
		if err != nil {
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}
	} else {
		refreshToken := ctx.GetCookieString("RefreshToken")
		if refreshToken == "" {
			ctx.WriteFailure(http.StatusUnauthorized, "login to your account")
			return
		}

		claims, err := m.jwt.Parse(refreshToken)
		if err != nil {
			ctx.WriteFailure(http.StatusBadRequest, "invalid refresh token")
			return
		}

		s, err := m.sessions.GetByUserId(context.Background(), claims.Id)
		if err != nil {
			ctx.WriteFailure(http.StatusBadRequest, "invalid refresh token")
			return
		}

		user, err := m.users.GetById(context.Background(), s.UserId)
		if err != nil {
			ctx.WriteFailure(http.StatusBadRequest, "invalid refresh token")
			return
		}

		newAccessToken, err := m.jwt.Generate(*user, m.tokenTTL)
		if err != nil {
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}

		newRefreshToken, err := m.jwt.Generate(*user, m.sessionTTL)
		if err != nil {
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}

		err = m.sessions.Update(context.Background(), newRefreshToken, s.Id, m.sessionTTL)
		if err != nil {
			ctx.WriteFailure(http.StatusInternalServerError, "server error")
			return
		}

		ctx.SetTokensCookie(newAccessToken, newRefreshToken, m.tokenTTL, m.sessionTTL)
	}
}
