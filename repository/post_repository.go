package repository

import (
	"context"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	GetPosts(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error)
	GetFollowingPosts(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error)
	GetPostById(ctx context.Context, id int) (*domain.Post, error)
	GetPostsByUserId(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetPostsCount(ctx context.Context) (int, error)
	GetFollowingPostsCount(ctx context.Context, userId int) (int, error)
	GetUserPostsCount(ctx context.Context, userId int) (int, error)
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPosts(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.SelectContext(ctx, &posts,
		`SELECT * FROM posts WHERE user_id != $1
         ORDER BY created_at DESC LIMIT $2 OFFSET $3 `,
		userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetFollowingPosts(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.SelectContext(ctx, &posts,
		`SELECT p.* FROM posts p
		 INNER JOIN followers f ON p.user_id = f.following_id
		 WHERE f.follower_id = $1
		 ORDER BY p.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPostById(ctx context.Context, id int) (*domain.Post, error) {
	post := domain.Post{}
	err := r.db.GetContext(ctx, &post, `SELECT * FROM posts WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostsByUserId(ctx context.Context, userId, limit, offset int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.SelectContext(ctx, &posts,
		`SELECT * FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO posts (user_id, ticker, body) VALUES ($1, $2, $3) RETURNING id`,
		post.UserId, post.Ticker, post.Body).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetPostById(ctx, id)
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = $1`, id)
	return err
}

func (r *postRepository) GetPostsCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM posts`)
	return count, err
}

func (r *postRepository) GetFollowingPostsCount(ctx context.Context, userId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM posts p
		 INNER JOIN followers f ON p.user_id = f.following_id
		 WHERE f.follower_id = $1`,
		userId)
	return count, err
}

func (r *postRepository) GetUserPostsCount(ctx context.Context, userId int) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM posts WHERE user_id = $1`, userId)
	return count, err
}
