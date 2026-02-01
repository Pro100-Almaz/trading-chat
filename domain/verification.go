package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrVerificationCodeExpired = errors.New("verification code has expired")
)

type VerificationCode struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required" example:"john@example.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required" example:"john@example.com"`
}

type VerificationResponse struct {
	Message string `json:"message" example:"Email verified successfully"`
}

type VerificationUseCase interface {
	VerifyEmail(ctx context.Context, email, code string) error
	ResendVerificationCode(ctx context.Context, email string) error
	SendVerificationCode(ctx context.Context, userId int, email string) error
}
