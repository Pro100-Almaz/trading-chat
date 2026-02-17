package domain

import "context"

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Successfully logged out"`
}

type LogoutUseCase interface {
	Logout(ctx context.Context, userId int, accessToken string, refreshToken string) error
}
