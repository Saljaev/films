package userhandler

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"tiny/internal/api/request"
	"tiny/internal/logger/sl"
	"tiny/internal/models"
)

func (h *UserHandler) Login(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "UserHandler - Login"

		log := log.With(
			slog.String("op", op),
		)

		var req request.UserRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))

			return
		}

		if req.Login == "" {
			log.Error("try to login without username")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("field login is required"))

			return
		}
		if req.Password == "" {
			log.Error("try to login without password")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("field password is required"))

			return
		}

		var user *models.User

		user, err = h.users.GetByLogin(r.Context(), req.Login)
		if err != nil {
			log.Error("failed to get user by login", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid login or password"))

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			log.Error("failed to compare passwords", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid login or password"))

			return
		}

		log.Info("login successful", slog.Any("id", user.Id))

		accessToken, err := h.jwt.Generate(*user)
		if err != nil {
			log.Error("failed to generate access token", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to create token"))

			return
		}

		refreshToken, err := h.jwt.Generate(*user)
		if err != nil {
			log.Error("failed to generate refresh token", sl.Err(err))

			w.Write([]byte("failed to create token"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		sessionCheck, err := h.sessions.GetByUserId(r.Context(), user.Id)

		if errors.Is(err, os.ErrNotExist) {
			_, err = h.sessions.Add(r.Context(), refreshToken, user.Id, h.sessionTTL)
			if err != nil {
				log.Error("failed to create session", sl.Err(err))

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("failed to create session"))

				return
			}
		} else if err != nil {
			log.Error("failed to get session", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))

			return
		} else {
			err = h.sessions.Update(r.Context(), refreshToken, sessionCheck.Id, h.sessionTTL)
			if err != nil {
				log.Error("failed to update session", sl.Err(err))

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal error"))

				return
			}
		}

		SetCookie(w, accessToken, refreshToken, h.tokenTTL, h.sessionTTL)

		w.Write([]byte("login successful"))
		w.WriteHeader(http.StatusOK)
	}
}

func SetCookie(w http.ResponseWriter, accessToken, refreshToken string, tokenTTL, sessionTTL time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Expires:  time.Now().Add(tokenTTL),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshToken,
		Expires: time.Now().Add(sessionTTL),
		Path:    "/",
	})
}