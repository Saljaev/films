package session

import (
	"context"
	"database/sql"
	"time"
	"tiny/internal/entities"
	"tiny/internal/usecase"
)

type UsersSessionRepo struct {
	*sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersSessionRepo {
	return &UsersSessionRepo{db}
}

var _ usecase.SessionRepo = (*UsersSessionRepo)(nil)

func (ur *UsersSessionRepo) Add(ctx context.Context, session entities.UserSession) (int, error) {
	panic("implement me")
}

func (ur *UsersSessionRepo) GetByUserId(ctx context.Context, userId int) (*entities.UserSession, error) {
	panic("implement me")
}

func (ur *UsersSessionRepo) Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error {
	panic("implement me")
}
