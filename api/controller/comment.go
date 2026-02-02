package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type CommentController struct {
	CommentUseCase domain.CommentUseCase
	Env            *bootstrap.Env
}

// GetComments godoc
// @Summary Get comments for a post
// @Description Get all comments for a post, paginated
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated comments"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id}/comments [get]
func (cc *CommentController) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	limit, offset := getPaginationParams(r)

	comments, err := cc.CommentUseCase.GetComments(r.Context(), postId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, comments)
}

// CreateComment godoc
// @Summary Add a comment to a post
// @Description Create a new comment on a post
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param request body domain.CreateCommentRequest true "Comment data"
// @Success 201 {object} domain.CommentResponse "Created comment"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id}/comments [post]
func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	var request domain.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := cc.CommentUseCase.CreateComment(r.Context(), userId, postId, &request)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusCreated, comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment (only owner can delete)
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /comments/{id} [delete]
func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	commentId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid comment id"})
		return
	}

	err = cc.CommentUseCase.DeleteComment(r.Context(), userId, commentId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}
