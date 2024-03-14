package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"tiny/internal/entities"
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
	const op = "UserUseCase - Register"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s - bcrypt.GenerateFromPassword: %w", op, err)
	}
	user := entities.User{
		Login:    u.Login,
		Password: string(hashPassword),
	}

	id, err := us.repo.Register(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%s - us.repo.Register: %w", op, err)
	}

	return id, nil
}

func (us *UsersUseCase) GetById(ctx context.Context, id int) (*models.User, error) {
	const op = "UsersUseCase - GetById"

	user, err := us.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s - us.repo.GetById: %w", op, err)
	}

	res := models.User{
		Id:       user.Id,
		Login:    user.Login,
		Role:     user.Role,
		Password: user.Password,
	}

	return &res, nil
}

func (us *UsersUseCase) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	const op = "UsersUseCase - GetByLogin"

	user, err := us.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("%s - us.repo.GetByLogin: %w", op, err)
	}

	res := models.User{
		Id:       user.Id,
		Login:    user.Login,
		Role:     user.Role,
		Password: user.Password,
	}

	return &res, nil
}

func (us *UsersUseCase) Delete(ctx context.Context, id int) error {
	const op = "UsersUseCase - Delete"

	err := us.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - us.repo.Delete: %w", op, err)
	}

	return nil
}
