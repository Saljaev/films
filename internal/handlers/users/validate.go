package userhandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"tiny/internal/logger/sl"
)

func (h *UserHandler) Validate(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "UserHandler - Validate"

			fmt.Println(log)

			log := log.With(
				slog.String("op", op),
			)

			token, err := r.Cookie("Token")

			if errors.Is(err, http.ErrNoCookie) {
				refreshToken, err := r.Cookie("RefreshToken")
				if err != nil {
					log.Info("user hase not token")

					w.Write([]byte("please, login to your account"))
					w.WriteHeader(http.StatusUnauthorized)

					return
				}

				claims, err := h.jwt.Parse(refreshToken.Value)
				if err != nil {
					log.Error("invalid refresh token token", sl.Err(err))

					w.Write([]byte("invalid refresh token"))
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				s, err := h.sessions.GetByUserId(r.Context(), claims.Id)
				if err != nil {
					log.Error("no user with this token", sl.Err(err))

					w.Write([]byte("invalid refresh token"))
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				user, err := h.users.GetById(r.Context(), s.UserId)

				newAccessToken, err := h.jwt.Generate(*user)
				if err != nil {
					log.Error("failed to create JWT token", sl.Err(err))

					w.Write([]byte("internal error"))
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
				newRefreshToken, err := h.jwt.Generate(*user)
				err = h.sessions.Update(r.Context(), newRefreshToken, s.Id, h.sessionTTL)
				if err != nil {
					log.Error("failed to update user's session", sl.Err(err))

					w.Write([]byte("internal error"))
					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				SetCookie(w, newAccessToken, newRefreshToken, h.tokenTTL, h.sessionTTL)

			} else {
				claims, err := h.jwt.Parse(token.Value)
				if err != nil {
					log.Error("invalid access token", sl.Err(err))

					w.Write([]byte("invalid access token"))
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				_, err = h.users.GetById(r.Context(), claims.Id)
				if err != nil {
					log.Error("failed to get user by login", sl.Err(err))

					w.Write([]byte("internal error"))
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
			}

			log.Info("user is verified")

			ctx := r.Context()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
