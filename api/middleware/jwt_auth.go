package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/internal/tokenutil"
	"github.com/Pro100-Almaz/trading-chat/repository"
	"github.com/Pro100-Almaz/trading-chat/utils"
)

func JwtAuthMiddleware(secret string, tokenBlacklistRepo repository.TokenBlacklistRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				t := strings.Split(authHeader, " ")
				if len(t) == 2 {
					authToken := t[1]

					// Check if token is blacklisted
					if tokenBlacklistRepo != nil {
						isBlacklisted, err := tokenBlacklistRepo.IsBlacklisted(r.Context(), authToken)
						if err != nil {
							utils.JSON(w, 401, domain.ErrorResponse{Message: "Failed to verify token"})
							return
						}
						if isBlacklisted {
							utils.JSON(w, 401, domain.ErrorResponse{Message: "Token has been revoked"})
							return
						}
					}

					authorized, err := tokenutil.IsAuthorized(authToken, secret)
					if err != nil {
						utils.JSON(w, 401, domain.ErrorResponse{Message: err.Error()})
						return
					}
					if authorized {
						userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
						if err != nil {
							utils.JSON(w, 401, domain.ErrorResponse{Message: err.Error()})
							return
						}
						// set user id to context
						ctx := context.WithValue(r.Context(), "user_id", userID)
						r = r.WithContext(ctx)
						next.ServeHTTP(w, r)
						return
					}
					utils.JSON(w, 401, domain.ErrorResponse{Message: domain.ErrUnauthorized.Error()})
					return
				}
				utils.JSON(w, 401, domain.ErrorResponse{Message: domain.ErrUnauthorized.Error()})
				return
			})
	}
}
