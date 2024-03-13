package models

import "time"

type UserSession struct {
	Id           int
	UserId       int
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExpiredAt    time.Time
}
