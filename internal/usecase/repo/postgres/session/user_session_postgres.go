package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
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
	const op = "UsersSessionRepo - Add"

	var id int

	query := "INSERT INTO users_sessions(user_id, refresh_token, expired_at, updated_at) " +
		"VALUES($1, $2, $3,$4)" +
		" RETURNING ID"
	err := ur.QueryRowContext(ctx, query, session.UserId, session.RefreshToken,
		session.ExpiredAt, session.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s - ur.QueryRowContext: %w", op, err)
	}

	return id, nil
}

func (ur *UsersSessionRepo) GetByUserId(ctx context.Context, userId int) (*entities.UserSession, error) {
	const op = "UsersSessionRepo - GetByUserId"

	session := entities.UserSession{}

	query := "SELECT * FROM users_sessions WHERE user_id = $1"
	row := ur.QueryRowContext(ctx, query, userId)
	err := row.Scan(&session.Id, &session.UserId, &session.RefreshToken, &session.ExpiredAt, &session.UpdatedAt, &session.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - us.QueryRowContext: %w", op, err)
	}

	return &session, nil
}

func (ur *UsersSessionRepo) Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error {
	const op = "UsersSessionRepo - Update"

	query := "UPDATE users_sessions SET refresh_token = $1, updated_at = $2 WHERE id = $3"
	_, err := ur.ExecContext(ctx, query, refreshToken, time.Now(), sessionId)
	if err != nil {
		return fmt.Errorf("%s - ur.QueryRowContext: %w", op, err)
	}

	return nil
}
