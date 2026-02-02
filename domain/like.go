package domain

import (
	"context"
	"time"
)

type Like struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	PostId    int       `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LikeUseCase interface {
	LikePost(ctx context.Context, userId, postId int) error
	UnlikePost(ctx context.Context, userId, postId int) error
}
