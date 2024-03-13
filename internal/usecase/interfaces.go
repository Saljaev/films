package usecase

import (
	"context"
	"time"
	"tiny/internal/entities"
	"tiny/internal/models"
)

type (
	Actors interface {
		Add(ctx context.Context, a models.Actor) (int, error)
		Update(ctx context.Context, a models.Actor) error
		Delete(ctx context.Context, id int) error
		GetAll(ctx context.Context) (map[*models.Actor][]*models.Films, error)
	}

	ActorsRepo interface {
		Add(ctx context.Context, a entities.Actors) (int, error)
		Update(ctx context.Context, a entities.Actors) error
		Delete(ctx context.Context, id int) error
		GetAll(ctx context.Context) (map[*models.Actor][]*models.Films, error)
	}

	Films interface {
		Add(ctx context.Context, f models.Films) (int, error)
		Update(ctx context.Context, f models.Films) error
		Delete(ctx context.Context, id int) error
		SearchByFragment(ctx context.Context, fragment, owner string) ([]*models.Films, error)
		RateByField(ctx context.Context, fragment string) ([]*models.Films, error)
	}

	FilmsRepo interface {
		Add(ctx context.Context, f entities.Films) (int, error)
		Update(ctx context.Context, f entities.Films) error
		Delete(ctx context.Context, id int) error
		SearchByFragment(ctx context.Context, fragment, owner string) ([]*entities.Films, error)
		RateByField(ctx context.Context, fragment string) ([]*entities.Films, error)
	}

	Users interface {
		Register(ctx context.Context, u models.User) (int, error)
		Login(ctx context.Context, u models.User) error
		GetById(ctx context.Context, id int) (*models.User, error)
		GetByLogin(ctx context.Context, login string) (*models.User, error)
		Delete(ctx context.Context, login string) error
	}

	UsersRepo interface {
		Register(ctx context.Context, u entities.User) (int, error)
		Login(ctx context.Context, u entities.User) error
		GetById(ctx context.Context, id int) (*entities.User, error)
		GetByLogin(ctx context.Context, login string) (*entities.User, error)
		Delete(ctx context.Context, login string) error
	}

	Sessions interface {
		Add(ctx context.Context, refreshToken string, userId int, sessonDuration time.Duration) (int, error)
		GetByUserId(ctx context.Context, userId int) (*models.UserSession, error)
		Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error
	}

	SessionRepo interface {
		Add(ctx context.Context, session entities.UserSession) (int, error)
		GetByUserId(ctx context.Context, userId int) (*entities.UserSession, error)
		Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error
	}
)
