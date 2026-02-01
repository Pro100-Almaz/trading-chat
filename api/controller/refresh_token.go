package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type RefreshTokenController struct {
	RefreshTokenUseCase domain.RefreshTokenUseCase
	Env                 *bootstrap.Env
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} domain.RefreshTokenResponse "Successfully refreshed tokens"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Router /refresh_token [post]
func (rtc *RefreshTokenController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request domain.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, refreshToken, err := rtc.RefreshTokenUseCase.RefreshToken(ctx, request, rtc.Env)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.JSON(w, http.StatusOK, refreshTokenResponse)
}
