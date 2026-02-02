package controller

import (
	"net/http"
	"strconv"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type FollowerController struct {
	FollowerUseCase domain.FollowerUseCase
	Env             *bootstrap.Env
}

// Follow godoc
// @Summary Follow a user
// @Description Follow another user
// @Tags Followers
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID to follow"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /users/{id}/follow [post]
func (fc *FollowerController) Follow(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	followingId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user id"})
		return
	}

	err = fc.FollowerUseCase.Follow(r.Context(), userId, followingId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user
// @Tags Followers
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID to unfollow"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /users/{id}/follow [delete]
func (fc *FollowerController) Unfollow(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	followingId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user id"})
		return
	}

	err = fc.FollowerUseCase.Unfollow(r.Context(), userId, followingId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}

// GetFollowers godoc
// @Summary Get user's followers
// @Description Get list of users who follow the specified user
// @Tags Followers
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated followers"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /users/{id}/followers [get]
func (fc *FollowerController) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user id"})
		return
	}

	limit, offset := getPaginationParams(r)

	followers, err := fc.FollowerUseCase.GetFollowers(r.Context(), userId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, followers)
}

// GetFollowing godoc
// @Summary Get user's following
// @Description Get list of users the specified user is following
// @Tags Followers
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated following"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /users/{id}/following [get]
func (fc *FollowerController) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user id"})
		return
	}

	limit, offset := getPaginationParams(r)

	following, err := fc.FollowerUseCase.GetFollowing(r.Context(), userId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, following)
}
