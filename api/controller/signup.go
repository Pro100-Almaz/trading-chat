package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
}

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account with email and password. A verification code will be sent to the email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.SignupRequest true "Signup credentials"
// @Success 201 {object} domain.SignupResponse "Account created, verification email sent"
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

	err := sc.SignupUseCase.SignUp(ctx, request)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	signupResponse := domain.SignupResponse{
		Message: "Account created successfully. Please check your email for verification code.",
		Email:   request.Email,
	}

	utils.JSON(w, http.StatusCreated, signupResponse)
}
