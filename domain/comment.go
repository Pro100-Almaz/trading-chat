package domain

import (
	"context"
	"time"
)

type Comment struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	PostId    int       `json:"post_id" db:"post_id"`
	Body      string    `json:"body" db:"body"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CommentResponse struct {
	Id        int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Author    Author    `json:"author"`
}

type CreateCommentRequest struct {
	Body string `json:"body" example:"Great analysis!"`
}

type CommentUseCase interface {
	GetComments(ctx context.Context, postId, limit, offset int) (*PaginatedResponse, error)
	CreateComment(ctx context.Context, userId, postId int, request *CreateCommentRequest) (*CommentResponse, error)
	DeleteComment(ctx context.Context, userId, commentId int) error
}
