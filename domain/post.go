package domain

import (
	"context"
	"time"
)

type Post struct {
	Id        int        `json:"id" db:"id"`
	UserId    int        `json:"user_id" db:"user_id"`
	Ticker    string     `json:"ticker" db:"ticker"`
	Body      string     `json:"body" db:"body"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type PostResponse struct {
	Id            int       `json:"id"`
	Ticker        string    `json:"ticker"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"created_at"`
	Author        Author    `json:"author"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	IsLiked       bool      `json:"is_liked"`
}

type Author struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	AvatarEmoji int    `json:"avatar_emoji"`
}

type CreatePostRequest struct {
	Ticker string `json:"ticker" example:"AAPL"`
	Body   string `json:"body" example:"I think this stock is going up!"`
}

type PostUseCase interface {
	GetGlobalFeed(ctx context.Context, userId, limit, offset int) (*PaginatedResponse, error)
	GetFollowingFeed(ctx context.Context, userId, limit, offset int) (*PaginatedResponse, error)
	GetUserPosts(ctx context.Context, currentUserId, targetUserId, limit, offset int) (*PaginatedResponse, error)
	GetPostById(ctx context.Context, userId, postId int) (*PostResponse, error)
	CreatePost(ctx context.Context, userId int, request *CreatePostRequest) (*PostResponse, error)
	DeletePost(ctx context.Context, userId, postId int) error
}
