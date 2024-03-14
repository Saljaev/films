package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny/internal/entities"
	"tiny/internal/models"
)

type UserSessionsUseCase struct {
	repo SessionRepo
}

var (
	ErrNoSession = errors.New("session not found")
)

var _ Sessions = (*UserSessionsUseCase)(nil)

func NewSessionUseCase(repo SessionRepo) *UserSessionsUseCase {
	return &UserSessionsUseCase{repo}
}

func (ss *UserSessionsUseCase) Add(ctx context.Context, refreshToken string, userId int, sessionDuration time.Duration) (int, error) {
	const op = "UserSessionsUseCase - Add"

	session := entities.UserSession{
		UserId:       userId,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(sessionDuration),
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	id, err := ss.repo.Add(ctx, session)
	if err != nil {
		return 0, fmt.Errorf("%s - uss.repo.Add: %w", op, err)
	}

	return id, nil
}

func (ss *UserSessionsUseCase) GetByUserId(ctx context.Context, userId int) (*models.UserSession, error) {
	const op = "UserSessionsUseCase - GetByUserId"

	session, err := ss.repo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s - ss.repo.GetByUserId: %w", op, err)
	}

	res := models.UserSession{
		Id:           session.Id,
		UserId:       userId,
		RefreshToken: session.RefreshToken,
		ExpiredAt:    session.ExpiredAt,
	}

	return &res, nil
}

func (ss UserSessionsUseCase) Update(ctx context.Context, refreshToken string, sessionId int, sessionDuration time.Duration) error {
	const op = "UserSessionsUseCase - Update"

	err := ss.repo.Update(ctx, refreshToken, sessionId, sessionDuration)
	if err != nil {
		return fmt.Errorf("%s - ss.repo.Update: %w", op, err)
	}

	return nil
}
