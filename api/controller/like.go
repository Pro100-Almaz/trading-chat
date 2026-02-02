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

type LikeController struct {
	LikeUseCase domain.LikeUseCase
	Env         *bootstrap.Env
}

// LikePost godoc
// @Summary Like a post
// @Description Like a post
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id}/like [post]
func (lc *LikeController) LikePost(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	err = lc.LikeUseCase.LikePost(r.Context(), userId, postId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}

// UnlikePost godoc
// @Summary Unlike a post
// @Description Unlike a post
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id}/like [delete]
func (lc *LikeController) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	err = lc.LikeUseCase.UnlikePost(r.Context(), userId, postId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}
