package repository

import (
	"context"
	"database/sql"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/jmoiron/sqlx"
)

type FollowerRepository interface {
	Follow(ctx context.Context, followerId, followingId int) error
	Unfollow(ctx context.Context, followerId, followingId int) error
	IsFollowing(ctx context.Context, followerId, followingId int) (bool, error)
	GetFollowers(ctx context.Context, userId, limit, offset int) ([]*domain.User, error)
	GetFollowing(ctx context.Context, userId, limit, offset int) ([]*domain.User, error)
	GetFollowersCount(ctx context.Context, userId int) (int, error)
	GetFollowingCount(ctx context.Context, userId int) (int, error)
}

type followerRepository struct {
	db *sqlx.DB
}

func NewFollowerRepository(db *sqlx.DB) FollowerRepository {
	return &followerRepository{db: db}
}

func (r *followerRepository) Follow(ctx context.Context, followerId, followingId int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO followers (follower_id, following_id) VALUES ($1, $2) ON CONFLICT (follower_id, following_id) DO NOTHING`,
		followerId, followingId)
	return err
}

func (r *followerRepository) Unfollow(ctx context.Context, followerId, followingId int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM followers WHERE follower_id = $1 AND following_id = $2`,
		followerId, followingId)
	return err
}

func (r *followerRepository) IsFollowing(ctx context.Context, followerId, followingId int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)`,
		followerId, followingId)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *followerRepository) GetFollowers(ctx context.Context, userId, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.* FROM users u
		 INNER JOIN followers f ON u.id = f.follower_id
		 WHERE f.following_id = $1
		 ORDER BY f.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *followerRepository) GetFollowing(ctx context.Context, userId, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.* FROM users u
		 INNER JOIN followers f ON u.id = f.following_id
		 WHERE f.follower_id = $1
		 ORDER BY f.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *followerRepository) GetFollowersCount(ctx context.Context, userId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM followers WHERE following_id = $1`,
		userId)
	return count, err
}

func (r *followerRepository) GetFollowingCount(ctx context.Context, userId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM followers WHERE follower_id = $1`,
		userId)
	return count, err
}
