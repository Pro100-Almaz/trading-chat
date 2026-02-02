package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type LikeRepository interface {
	LikePost(ctx context.Context, userId, postId int) error
	UnlikePost(ctx context.Context, userId, postId int) error
	IsLiked(ctx context.Context, userId, postId int) (bool, error)
	GetLikesCount(ctx context.Context, postId int) (int, error)
}

type likeRepository struct {
	db *sqlx.DB
}

func NewLikeRepository(db *sqlx.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) LikePost(ctx context.Context, userId, postId int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT (user_id, post_id) DO NOTHING`,
		userId, postId)
	return err
}

func (r *likeRepository) UnlikePost(ctx context.Context, userId, postId int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM likes WHERE user_id = $1 AND post_id = $2`,
		userId, postId)
	return err
}

func (r *likeRepository) IsLiked(ctx context.Context, userId, postId int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`,
		userId, postId)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *likeRepository) GetLikesCount(ctx context.Context, postId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM likes WHERE post_id = $1`,
		postId)
	return count, err
}
