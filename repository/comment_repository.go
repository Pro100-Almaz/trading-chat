package repository

import (
	"context"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	GetCommentsByPostId(ctx context.Context, postId, limit, offset int) ([]*domain.Comment, error)
	GetCommentById(ctx context.Context, id int) (*domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	GetCommentsCount(ctx context.Context, postId int) (int, error)
}

type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) GetCommentsByPostId(ctx context.Context, postId, limit, offset int) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	err := r.db.SelectContext(ctx, &comments,
		`SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		postId, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) GetCommentById(ctx context.Context, id int) (*domain.Comment, error) {
	comment := domain.Comment{}
	err := r.db.GetContext(ctx, &comment, `SELECT * FROM comments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO comments (user_id, post_id, body) VALUES ($1, $2, $3) RETURNING id`,
		comment.UserId, comment.PostId, comment.Body).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetCommentById(ctx, id)
}

func (r *commentRepository) DeleteComment(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)
	return err
}

func (r *commentRepository) GetCommentsCount(ctx context.Context, postId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM comments WHERE post_id = $1`, postId)
	return count, err
}
