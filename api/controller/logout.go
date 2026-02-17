package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type LogoutController struct {
	LogoutUseCase domain.LogoutUseCase
}

// Logout godoc
// @Summary Logout user
// @Description Logout the currently authenticated user by blacklisting their tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.LogoutRequest false "Optional refresh token to also blacklist"
// @Success 200 {object} domain.LogoutResponse "Successfully logged out"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /logout [post]
func (lc *LogoutController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (set by JWT middleware)
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		userIDStr, ok := ctx.Value("user_id").(string)
		if !ok {
			log.Error("User ID not found in context")
			utils.JSON(w, http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
			return
		}
		var err error
		userID, err = strconv.Atoi(userIDStr)
		if err != nil {
			log.Error("Invalid user ID in context: ", err)
			utils.JSON(w, http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
			return
		}
	}

	// Extract access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Error("Authorization header is missing")
		utils.JSON(w, http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	// Remove "Bearer " prefix
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	if accessToken == authHeader {
		log.Error("Invalid authorization header format")
		utils.JSON(w, http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid authorization header"})
		return
	}

	// Get optional refresh token from request body
	var request domain.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// It's okay if body is empty or invalid, we'll just logout with access token only
		log.Debug("No refresh token provided in logout request")
	}

	// Logout
	err := lc.LogoutUseCase.Logout(ctx, userID, accessToken, request.RefreshToken)
	if err != nil {
		log.Error("Failed to logout: ", err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, domain.LogoutResponse{
		Message: "Successfully logged out",
	})
}
