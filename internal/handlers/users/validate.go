package userhandler

import (
	"errors"
	"net/http"
)

func (h *UserHandler) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "UserHandler - Validate"

		token, err := r.Cookie("Token")

		if errors.Is(err, http.ErrNoCookie) {
			refreshToken, err := r.Cookie("RefreshToken")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("please, login to your account"))

				return
			}

			claims, err := h.jwt.Parse(refreshToken.Value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid refresh token"))

				return
			}

			s, err := h.sessions.GetByUserId(r.Context(), claims.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid refresh token"))

				return
			}

			user, err := h.users.GetById(r.Context(), s.UserId)

			newAccessToken, err := h.jwt.Generate(*user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal error"))

				return
			}
			newRefreshToken, err := h.jwt.Generate(*user)
			err = h.sessions.Update(r.Context(), newRefreshToken, s.Id, h.sessionTTL)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal error"))

				return
			}

			SetCookie(w, newAccessToken, newRefreshToken, h.tokenTTL, h.sessionTTL)

		} else {
			claims, err := h.jwt.Parse(token.Value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid access token"))

				return
			}

			_, err = h.users.GetById(r.Context(), claims.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal error"))

				return
			}
		}

		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
