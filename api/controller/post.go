package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type PostController struct {
	PostUseCase domain.PostUseCase
	Env         *bootstrap.Env
}

// GetGlobalFeed godoc
// @Summary Get global feed
// @Description Get all posts, paginated
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated posts"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts [get]
func (pc *PostController) GetGlobalFeed(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	limit, offset := getPaginationParams(r)

	feed, err := pc.PostUseCase.GetGlobalFeed(r.Context(), userId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, feed)
}

// GetFollowingFeed godoc
// @Summary Get following feed
// @Description Get posts from followed users, paginated
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated posts"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/following [get]
func (pc *PostController) GetFollowingFeed(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	limit, offset := getPaginationParams(r)

	feed, err := pc.PostUseCase.GetFollowingFeed(r.Context(), userId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, feed)
}

// GetPost godoc
// @Summary Get a single post
// @Description Get a post by ID
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} domain.PostResponse "Post details"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id} [get]
func (pc *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	post, err := pc.PostUseCase.GetPostById(r.Context(), userId, postId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, post)
}

// GetUserPosts godoc
// @Summary Get user's posts
// @Description Get posts by a specific user
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} domain.PaginatedResponse "Paginated posts"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/user/{id} [get]
func (pc *PostController) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	targetUserId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user id"})
		return
	}

	limit, offset := getPaginationParams(r)

	posts, err := pc.PostUseCase.GetUserPosts(r.Context(), userId, targetUserId, limit, offset)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, posts)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreatePostRequest true "Post data"
// @Success 201 {object} domain.PostResponse "Created post"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts [post]
func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	var request domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	post, err := pc.PostUseCase.CreatePost(r.Context(), userId, &request)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusCreated, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post (only owner can delete)
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {string} string "Success"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /posts/{id} [delete]
func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r)

	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: "invalid post id"})
		return
	}

	err = pc.PostUseCase.DeletePost(r.Context(), userId, postId)
	if err != nil {
		log.Error(err)
		utils.JSON(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	utils.JSON(w, http.StatusOK, "Success")
}

func getUserIdFromContext(r *http.Request) int {
	id := fmt.Sprintf("%v", r.Context().Value("user_id"))
	userId, _ := strconv.Atoi(id)
	return userId
}

func getPaginationParams(r *http.Request) (limit, offset int) {
	limit = 10
	offset = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}
