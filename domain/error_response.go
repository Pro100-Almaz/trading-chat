package domain

import "errors"

type ErrorResponse struct {
	Message string `json:"message"`
}

// Error List:
var (
	ErrUserNotAllowed            = errors.New("user not allowed")
	ErrUserNotFound              = errors.New("user not found")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrUserShouldLoginWithGoogle = errors.New("user should login with Google")
	ErrCodeExchangeWrong         = errors.New("code exchange wrong")
	ErrFailedGetGoogleUser       = errors.New("failed to get google user")
	ErrFailedToReadResponse      = errors.New("failed to read response")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrInvalidToken              = errors.New("invalid token")
	ErrEmailNotVerified          = errors.New("email not verified, please verify your email first")
	ErrInvalidVerificationCode   = errors.New("invalid or expired verification code")
	ErrUserAlreadyVerified       = errors.New("user is already verified")
	ErrFailedToSendEmail         = errors.New("failed to send verification email")
)
