package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
	Env           *bootstrap.Env
}

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.SignupRequest true "Signup credentials"
// @Success 200 {object} domain.SignupResponse "Successfully registered"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Router /signup [post]
func (sc *SignupController) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request domain.SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, refreshToken, err := sc.SignupUseCase.SignUp(ctx, request, sc.Env)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	signupResponse := domain.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.JSON(w, http.StatusOK, signupResponse)
}
