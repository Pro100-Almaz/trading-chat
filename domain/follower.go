package domain

import (
	"context"
	"time"
)

type Follower struct {
	Id          int       `json:"id" db:"id"`
	FollowerId  int       `json:"follower_id" db:"follower_id"`
	FollowingId int       `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type FollowUserResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	AvatarEmoji int    `json:"avatar_emoji"`
}

type FollowerUseCase interface {
	Follow(ctx context.Context, followerId, followingId int) error
	Unfollow(ctx context.Context, followerId, followingId int) error
	GetFollowers(ctx context.Context, userId, limit, offset int) (*PaginatedResponse, error)
	GetFollowing(ctx context.Context, userId, limit, offset int) (*PaginatedResponse, error)
	IsFollowing(ctx context.Context, followerId, followingId int) (bool, error)
}
