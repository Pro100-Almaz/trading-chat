package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type VerificationController struct {
	VerificationUseCase domain.VerificationUseCase
}

// VerifyEmail godoc
// @Summary Verify email with code
// @Description Verify user's email address using the 6-digit code sent to their email
// @Tags Verification
// @Accept json
// @Produce json
// @Param request body domain.VerifyEmailRequest true "Email and verification code"
// @Success 200 {object} domain.VerificationResponse "Email verified successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad request or invalid code"
// @Router /verify-email [post]
func (vc *VerificationController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request domain.VerifyEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := vc.VerificationUseCase.VerifyEmail(ctx, request.Email, request.Code)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, domain.VerificationResponse{Message: "Email verified successfully"})
}

// ResendVerificationCode godoc
// @Summary Resend verification code
// @Description Resend a new verification code to the user's email
// @Tags Verification
// @Accept json
// @Produce json
// @Param request body domain.ResendVerificationRequest true "User email"
// @Success 200 {object} domain.VerificationResponse "Verification code sent"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Router /resend-verification [post]
func (vc *VerificationController) ResendVerificationCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request domain.ResendVerificationRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := vc.VerificationUseCase.ResendVerificationCode(ctx, request.Email)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, domain.VerificationResponse{Message: "Verification code sent to your email"})
}
