package users

import (
	"context"
	"database/sql"
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
	panic("implement me")
}

func (us *UsersRepo) Login(ctx context.Context, u entities.User) error {
	panic("implement me")
}

func (us *UsersRepo) GetById(ctx context.Context, id int) (*entities.User, error) {
	panic("implement me")
}

func (us *UsersRepo) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	panic("implement me")
}

func (us *UsersRepo) Delete(ctx context.Context, login string) error {
	panic("implement me")
}
