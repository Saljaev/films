package usecase

import (
	"context"
	"tiny/internal/models"
)

type UsersUseCase struct {
	repo UsersRepo
}

var _ Users = (*UsersUseCase)(nil)

func NewUsersUseCase(repo UsersRepo) *UsersUseCase {
	return &UsersUseCase{repo}
}

func (us *UsersUseCase) Register(ctx context.Context, u models.User) (int, error) {
	panic("implement me")
}

func (us *UsersUseCase) Login(ctx context.Context, u models.User) error {
	panic("implement me")
}

func (us *UsersUseCase) GetById(ctx context.Context, id int) (*models.User, error) {
	panic("implement me")
}

func (us *UsersUseCase) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	panic("implement me")
}

func (us *UsersUseCase) Delete(ctx context.Context, login string) error {
	panic("implement me")
}
