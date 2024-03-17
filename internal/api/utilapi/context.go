package utilapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"tiny/internal/logger/sl"
)

type Error struct {
	ErrorMessage string `json:"error_message"`
}

type APIContext struct {
	w             http.ResponseWriter
	r             *http.Request
	log           *slog.Logger
	next          HandlerFunc
	writeResponse bool
}

func NewAPIContext(
	log *slog.Logger,
) *APIContext {
	return &APIContext{
		log:           log,
		writeResponse: false,
	}
}

func (ctx *APIContext) Error(msg string, err error) {
	ctx.log.Error(msg, sl.Err(err))
}

func (ctx *APIContext) Info(msg string, key string, value interface{}) {
	ctx.log.Info(msg, slog.Any(key, value))
}

type validator interface {
	IsValid() bool
}

func (ctx *APIContext) Decode(dest validator) error {
	err := json.NewDecoder(ctx.r.Body).Decode(&dest)
	if err != nil || !dest.IsValid() {
		if err == nil {
			err = errors.New("invalid request")
		}
		ctx.Error("error", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid request")

		return err
	}

	return nil
}

func (ctx *APIContext) WriteFailure(code int, msg string) {
	ctx.w.WriteHeader(code)

	data, _ := json.Marshal(Error{ErrorMessage: msg})

	_, err := ctx.w.Write(data)
	if err != nil {
		ctx.Error("response error", err)
	}
	ctx.writeResponse = true
	ctx.r.Context().Done()
}

func (ctx *APIContext) SuccessWithData(data interface{}) {
	jsonData, _ := json.Marshal(data)

	ctx.w.WriteHeader(http.StatusOK)
	_, _ = ctx.w.Write(jsonData)
}

func (ctx *APIContext) SetTokensCookie(accessToken, refreshToken string, tokenTTL, sessionTTL time.Duration) {
	http.SetCookie(ctx.w, &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Expires:  time.Now().Add(tokenTTL),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(ctx.w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
		Path:     "/",
	})
}

func (ctx *APIContext) GetURLParam(name string) string {
	return ctx.r.URL.Query().Get(name)
}

func (ctx *APIContext) GetCookieString(name string) string {
	cookie, err := ctx.r.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (ctx *APIContext) ContextDone() {
	ctx.r.Context().Done()
}
