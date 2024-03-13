package usecase

import (
	"context"
	"errors"
	"time"
	"tiny/internal/models"
)

type UserSessionUseCase struct {
	repo SessionRepo
}

var (
	ErrNoSession = errors.New("session not found")
)

var _ Sessions = (*UserSessionUseCase)(nil)

func NewSessionUseCase(repo SessionRepo) *UserSessionUseCase {
	return &UserSessionUseCase{repo}
}

func (ss *UserSessionUseCase) Add(ctx context.Context, refreshToken string, userId int, sessionDuration time.Duration) (int, error) {
	panic("implement me")
}

func (ss *UserSessionUseCase) GetByUserId(ctx context.Context, userId int) (*models.UserSession, error) {
	panic("implement me")
}

func (ss UserSessionUseCase) Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error {
	panic("implement me")
}
