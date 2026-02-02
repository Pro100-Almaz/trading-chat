package domain

import "context"

type SignupRequest struct {
	Name        string `json:"name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john@example.com"`
	Password    string `json:"password" binding:"required" example:"password123"`
	AvatarEmoji *int   `json:"avatar_emoji,omitempty" example:"5"`
}

type SignupResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type SignupUseCase interface {
	SignUp(ctx context.Context, request SignupRequest) error
}
