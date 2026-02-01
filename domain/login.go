package domain

import (
	"context"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUseCase interface {
	Login(ctx context.Context, request LoginRequest, env *bootstrap.Env) (accessToken string, refreshToken string, err error)
}
