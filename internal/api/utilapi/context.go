package utilapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	userhandler "tiny/internal/api/handlers/users"
	"tiny/internal/api/response"
	"tiny/internal/logger/sl"
	"tiny/internal/usecase"
)

type APIContext struct {
	w          http.ResponseWriter
	r          *http.Request
	users      usecase.Users
	sessions   usecase.Sessions
	jwt        *userhandler.JWTManager
	log        *slog.Logger
	tokenTTL   time.Duration
	sessionTTL time.Duration
	next       HandlerFunc
}

func NewAPIContext(
	users usecase.Users,
	sessions usecase.Sessions,
	jwt *userhandler.JWTManager,
	log *slog.Logger,
	tokenTTL time.Duration,
	sessionTTL time.Duration,
) *APIContext {
	return &APIContext{
		users:      users,
		sessions:   sessions,
		jwt:        jwt,
		log:        log,
		tokenTTL:   tokenTTL,
		sessionTTL: sessionTTL,
	}
}

func (ctx *APIContext) Error(msg string, args ...interface{}) {
	ctx.log.Error(msg, args)
}

func (ctx *APIContext) Info(msg string, args ...interface{}) {
	ctx.log.Info(msg, args)
}

type validator interface {
	IsValid() bool
}

func (ctx *APIContext) Decode(dest validator) {
	err := json.NewDecoder(ctx.r.Body).Decode(&dest)
	if err != nil || !dest.IsValid() {
		ctx.Error("error", sl.Err(err))
		ctx.WriteFailure(http.StatusBadRequest, "invalid request")
	}
}

func (ctx *APIContext) WriteFailure(code int, msg string) {
	ctx.w.WriteHeader(code)

	data, _ := json.Marshal(response.Error{ErrorMessage: msg})

	_, err := ctx.w.Write(data)
	if err != nil {
		ctx.Error("response error", sl.Err(err))
	}

	ctx.r.Context().Done()
}

func (ctx *APIContext) SuccessWithData(data interface{}) {
	jsonData, _ := json.Marshal(data)

	ctx.w.WriteHeader(http.StatusOK)
	_, _ = ctx.w.Write(jsonData)
}

func (ctx *APIContext) SetTokensCookie(accessToken, refreshToken string) {
	http.SetCookie(ctx.w, &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Expires:  time.Now().Add(ctx.tokenTTL),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(ctx.w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(ctx.sessionTTL),
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
