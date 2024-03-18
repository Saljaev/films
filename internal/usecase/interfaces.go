package usecase

import (
	"context"
	"time"
	"tiny/internal/entities"
	"tiny/internal/models"
)

type (
	Actors interface {
		Add(ctx context.Context, a models.Actors) (int, error)
		GetById(ctx context.Context, id int) (*models.Actors, error)
		Update(ctx context.Context, a models.Actors) error
		Delete(ctx context.Context, id int) error
		GetAll(ctx context.Context) ([]*models.Actors, error)
	}

	ActorsRepo interface {
		Add(ctx context.Context, a entities.Actors) (int, error)
		GetById(ctx context.Context, id int) (*entities.Actors, error)
		Update(ctx context.Context, a entities.Actors) error
		Delete(ctx context.Context, id int) error
		GetAll(ctx context.Context) ([]*entities.Actors, error)
	}

	Films interface {
		GetById(ctx context.Context, id int) (*models.Films, error)
		Add(ctx context.Context, f models.Films) (int, error)
		Update(ctx context.Context, f models.Films) error
		Delete(ctx context.Context, id int) error
		SearchByFilmName(ctx context.Context, name string) ([]*models.Films, error)
		SearchByActorName(ctx context.Context, firstName, lastName string) ([]*models.Films, error)
		RateByField(ctx context.Context, fragment string, increasing bool) ([]*models.Films, error)
		DeleteActor(ctx context.Context, filmID, actorID int) error
	}

	FilmsRepo interface {
		GetById(ctx context.Context, id int) (*entities.Films, error)
		Add(ctx context.Context, f entities.Films) (int, error)
		Update(ctx context.Context, f entities.Films) error
		Delete(ctx context.Context, id int) error
		SearchByFilmName(ctx context.Context, name string) ([]*entities.Films, error)
		SearchByActorName(ctx context.Context, firstName, lastName string) ([]*entities.Films, error)
		RateByField(ctx context.Context, fragment string, increasing string) ([]*entities.Films, error)
		DeleteActor(ctx context.Context, filmID, actorID int) error
	}

	Users interface {
		Register(ctx context.Context, u models.User) (int, error)
		GetById(ctx context.Context, id int) (*models.User, error)
		GetByLogin(ctx context.Context, login string) (*models.User, error)
		Delete(ctx context.Context, id int) error
	}

	UsersRepo interface {
		Register(ctx context.Context, u entities.User) (int, error)
		GetById(ctx context.Context, id int) (*entities.User, error)
		GetByLogin(ctx context.Context, login string) (*entities.User, error)
		Delete(ctx context.Context, id int) error
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
