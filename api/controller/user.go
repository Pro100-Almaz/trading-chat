package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	log "github.com/sirupsen/logrus"
)

type UserController struct {
	UserUseCase domain.UserUseCase
	Env         *bootstrap.Env
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.UserResponse "List of users"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /user/all [get]
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "user_id", r.Context().Value("user_id"))

	users, err := uc.UserUseCase.GetUsers(ctx)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, users)
}

// GetUserById godoc
// @Summary Get current user
// @Description Get the profile of the currently authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserResponse "User profile"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /user [get]
func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "user_id", r.Context().Value("user_id"))

	id := fmt.Sprintf("%v", ctx.Value("user_id"))

	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := uc.UserUseCase.GetUserById(ctx, intId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update current user
// @Description Update the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UserUpdateRequest true "User data to update"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /user [put]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "user_id", r.Context().Value("user_id"))

	var user *domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	id := fmt.Sprintf("%v", ctx.Value("user_id"))

	userId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user.Id = userId

	err = uc.UserUseCase.UpdateUser(ctx, user)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}

// DeleteUser godoc
// @Summary Delete current user
// @Description Delete the account of the currently authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /user [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "user_id", r.Context().Value("user_id"))

	id, err := strconv.Atoi(fmt.Sprintf("%v", ctx.Value("user_id")))
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err = uc.UserUseCase.DeleteUser(ctx, id)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}
