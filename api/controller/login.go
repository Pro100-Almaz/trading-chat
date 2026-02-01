package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type LoginController struct {
	LoginUseCase domain.LoginUseCase
	Env          *bootstrap.Env
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse "Successfully logged in"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Router /login [post]
func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var request domain.LoginRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, refreshToken, err := lc.LoginUseCase.Login(ctx, request, lc.Env)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	response := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.JSON(w, http.StatusOK, response)
}
