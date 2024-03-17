package utilapi

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	log *slog.Logger
}

func NewRouter(log *slog.Logger) *Router {
	return &Router{
		mux: http.NewServeMux(),
		log: log,
	}
}

func (router *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := NewAPIContext(router.log)

		ctx.w = w
		ctx.r = r

		for _, h := range handlerFuncs {
			if !ctx.writeResponse {
				ctx.log = slog.With(slog.String("op", fmt.Sprintf("%s", pattern)))
				h(ctx)
			}
		}
	})
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
