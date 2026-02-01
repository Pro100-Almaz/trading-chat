package domain

import (
	"context"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
)

type SignupRequest struct {
	Name        string `json:"name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john@example.com"`
	Password    string `json:"password" binding:"required" example:"password123"`
	AvatarEmoji *int   `json:"avatar_emoji,omitempty" example:"5"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupUseCase interface {
	SignUp(ctx context.Context, request SignupRequest, env *bootstrap.Env) (accessToken string, refreshToken string, err error)
}
