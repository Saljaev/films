package userhandler

import (
	"net/http"
)

func (h *UserHandler) CheckRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "UserHandler - Validate"

		token, err := r.Cookie("Token")
		claims, err := h.jwt.Parse(token.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))

			return
		}

		user, err := h.users.GetById(r.Context(), claims.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))

			return
		}

		if user.Role == "user" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("not enough rights"))

			return
		}

		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
