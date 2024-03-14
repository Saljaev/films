package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"tiny/internal/entities"
	"tiny/internal/usecase"
)

type UsersRepo struct {
	*sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db}
}

var _ usecase.UsersRepo = (*UsersRepo)(nil)

func (us *UsersRepo) Register(ctx context.Context, u entities.User) (int, error) {
	const op = "UsersRepo - Register"

	var id int

	query := "INSERT INTO users(login, password) VALUES($1, $2) RETURNING ID"
	err := us.QueryRowContext(ctx, query, u.Login, u.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s - us.QueryRowContext: %w", op, err)
	}

	return id, nil
}

func (us *UsersRepo) GetById(ctx context.Context, id int) (*entities.User, error) {
	const op = "UsersRepo - GetById"

	user := entities.User{}

	query := "SELECT * FROM users WHERE id = $1"
	row := us.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.Login, &user.Role, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - us.QueryRowContext: %w", op, err)
	}

	return &user, nil
}

func (us *UsersRepo) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	const op = "UsersRepo - GetByLogin"

	user := entities.User{}

	query := "SELECT * FROM users WHERE login = $1"
	row := us.QueryRowContext(ctx, query, login)
	err := row.Scan(&user.Id, &user.Login, &user.Role, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - us.QueryRowContext: %w", op, err)
	}

	return &user, nil
}

func (us *UsersRepo) Delete(ctx context.Context, id int) error {
	const op = "UserRepo - Delete"

	tx, err := us.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s - us.BeginTx: %w", op, err)
	}

	defer tx.Rollback()

	query := "DELETE FROM users_sessions WHERE user_id = $1"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s - tx.ExecContext: %w", op, err)
	}

	query = "DELETE FROM users WHERE id = $1"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s - tx.ExecContext: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return nil
}
