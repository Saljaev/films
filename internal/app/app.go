package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	actorshandler "tiny/internal/api/handlers/actors"
	filmshandler "tiny/internal/api/handlers/films"
	"tiny/internal/api/handlers/middleware"
	"tiny/internal/api/utilapi"
	"tiny/internal/config"
	"tiny/internal/logger/slogpretty"
	actorrepo "tiny/internal/usecase/repo/postgres/actors"
	filmrepo "tiny/internal/usecase/repo/postgres/films"
	sessionrepo "tiny/internal/usecase/repo/postgres/session"
	userrepo "tiny/internal/usecase/repo/postgres/users"
	"tiny/pkg"

	_ "github.com/lib/pq"
	userhandler "tiny/internal/api/handlers/users"
	"tiny/internal/logger/sl"
	"tiny/internal/usecase"
)

func Run() {
	cfg := config.ConfgiLoad()

	log := setupPrettySlog()

	// TODO: init db
	PG_URL := cfg.StoragePath

	db, err := sql.Open("postgres", PG_URL)
	if err != nil {
		log.Error("failed to connect db")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Error("failed to check db connection")
	} else {
		log.Info("successful connect to db")
	}

	users := usecase.NewUsersUseCase(userrepo.NewUsersRepo(db))
	sessions := usecase.NewSessionUseCase(sessionrepo.NewUsersRepo(db))
	actors := usecase.NewActorsUseCase(actorrepo.NewActorsRepo(db))
	films := usecase.NewFilmsUseCase(filmrepo.NewFilmsRepo(db))

	jwtManager := pkg.NewJWTManager(cfg.Secret, cfg.TokenTTL)

	userHandler := userhandler.NewUserHandler(users, sessions, *jwtManager, cfg.TokenTTL, cfg.SessionTTL)
	filmsHandler := filmshandler.NewFilmsHandler(films)
	actorsHandler := actorshandler.NewActorsHandler(actors)

	m := middleware.NewMiddleware(users, sessions, *jwtManager, cfg.TokenTTL, cfg.SessionTTL)

	router := utilapi.NewRouter(log)

	router.Handle("POST /user/register", userHandler.Register)
	router.Handle("POST /user/login", userHandler.Login)
	router.Handle("DELETE /user/", m.Validate, m.CheckRole, userHandler.Delete)

	router.Handle("POST /actor", m.Validate, m.CheckRole, actorsHandler.Create)
	router.Handle("POST /actor/", m.Validate, m.CheckRole, actorsHandler.Update)
	router.Handle("DELETE /actor/", m.Validate, m.CheckRole, actorsHandler.Delete)
	router.Handle("GET /actors", actorsHandler.GetWithFilms)

	router.Handle("POST /film", m.Validate, m.CheckRole, filmsHandler.Add)
	router.Handle("UPDATE /film/", m.Validate, m.CheckRole, filmsHandler.Update)
	router.Handle("DELETE /film/", m.Validate, m.CheckRole, filmsHandler.Delete)
	router.Handle("GET /film.search_by_fragment", filmsHandler.SearchByFragment)
	router.Handle("GET /film.get", filmsHandler.GetWithSort)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		srv.ListenAndServe()
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
