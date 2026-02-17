package usecase

import (
	"context"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"
	"github.com/golang-jwt/jwt/v4"

	log "github.com/sirupsen/logrus"
)

type logoutUseCase struct {
	tokenBlacklistRepository repository.TokenBlacklistRepository
	contextTimeout           time.Duration
}

func NewLogoutUseCase(
	tokenBlacklistRepo repository.TokenBlacklistRepository,
	timeout time.Duration,
) domain.LogoutUseCase {
	return &logoutUseCase{
		tokenBlacklistRepository: tokenBlacklistRepo,
		contextTimeout:           timeout,
	}
}

func (lu *logoutUseCase) Logout(ctx context.Context, userId int, accessToken string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	// Extract expiration from access token
	accessTokenExp, err := extractTokenExpiration(accessToken)
	if err != nil {
		log.Error("Failed to extract access token expiration: ", err)
		// Continue anyway, use a default expiration
		accessTokenExp = time.Now().Add(2 * time.Hour)
	}

	// Blacklist access token
	err = lu.tokenBlacklistRepository.AddToBlacklist(ctx, accessToken, accessTokenExp)
	if err != nil {
		log.Error("Failed to blacklist access token: ", err)
		return err
	}

	// Blacklist refresh token if provided
	if refreshToken != "" {
		refreshTokenExp, err := extractTokenExpiration(refreshToken)
		if err != nil {
			log.Error("Failed to extract refresh token expiration: ", err)
			// Continue anyway, use a default expiration
			refreshTokenExp = time.Now().Add(168 * time.Hour)
		}

		err = lu.tokenBlacklistRepository.AddToBlacklist(ctx, refreshToken, refreshTokenExp)
		if err != nil {
			log.Error("Failed to blacklist refresh token: ", err)
			return err
		}
	}

	log.Infof("User %d logged out successfully", userId)
	return nil
}

// extractTokenExpiration extracts the expiration time from a JWT token
func extractTokenExpiration(tokenString string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			return time.Unix(int64(exp), 0), nil
		}
	}

	return time.Time{}, domain.ErrInvalidToken
}
