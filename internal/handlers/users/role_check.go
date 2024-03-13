package userhandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"tiny/internal/logger/sl"
)

func (h *UserHandler) CheckRole(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "UserHandler - Validate"

			fmt.Println(log)

			log := log.With(
				slog.String("op", op),
			)

			token, err := r.Cookie("Token")
			claims, err := h.jwt.Parse(token.Value)
			if err != nil {
				log.Error("failed to parse access token", sl.Err(err))

				w.Write([]byte("internal error"))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			user, err := h.users.GetById(r.Context(), claims.Id)
			if err != nil {
				log.Error("failed to get user", sl.Err(err))

				w.Write([]byte("internal error"))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_ = user
			panic("implement me")

			log.Info("user is verified")

			ctx := r.Context()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
