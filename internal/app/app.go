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
	"tiny/internal/config"
	actorshandler "tiny/internal/handlers/actors"
	filmshandler "tiny/internal/handlers/films"
	"tiny/internal/logger/slogpretty"
	actorrepo "tiny/internal/usecase/repo/postgres/actors"
	filmrepo "tiny/internal/usecase/repo/postgres/films"
	sessionrepo "tiny/internal/usecase/repo/postgres/session"
	userrepo "tiny/internal/usecase/repo/postgres/users"

	_ "github.com/lib/pq"
	userhandler "tiny/internal/handlers/users"
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

	jwtManager := userhandler.NewJWTManager(cfg.Secret, cfg.TokenTTL)

	userHandler := userhandler.NewUserHandler(users, sessions, *jwtManager, cfg.TokenTTL, cfg.SessionTTL)
	actorsHandler := actorshandler.NewActorsHandler(actors)
	filmsHandler := filmshandler.NewFilmsHandler(films)

	// TODO: init router
	router := http.NewServeMux()

	router.Handle("/user.delete/", userHandler.Validate(userHandler.CheckRole(userHandler.Delete(log))))
	router.Handle("/user.create", userHandler.Register(log))
	router.Handle("/user.login", userHandler.Login(log))

	router.Handle("POST /actor.create/", actorsHandler.Add(log))
	router.Handle("POST /actor.update/", actorsHandler.Update(log))

	router.Handle("/film.add", filmsHandler.Add(log))
	router.Handle("/film.update/", filmsHandler.Update(log))

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
