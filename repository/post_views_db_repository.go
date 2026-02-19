package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostViewsDBRepository interface {
	UpsertViewCounts(ctx context.Context, viewCounts map[int]int64) error
	GetViewCount(ctx context.Context, postId int) (int64, error)
	GetMultipleViewCounts(ctx context.Context, postIds []int) (map[int]int64, error)
}

type postViewsDBRepository struct {
	db *sqlx.DB
}

func NewPostViewsDBRepository(db *sqlx.DB) PostViewsDBRepository {
	return &postViewsDBRepository{
		db: db,
	}
}

// UpsertViewCounts updates or inserts view counts in batch
func (r *postViewsDBRepository) UpsertViewCounts(ctx context.Context, viewCounts map[int]int64) error {
	if len(viewCounts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// PostgreSQL INSERT ... ON CONFLICT for batch upsert
	query := `
		INSERT INTO post_views (post_id, views_count, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id)
		DO UPDATE SET
			views_count = post_views.views_count + EXCLUDED.views_count,
			updated_at = EXCLUDED.updated_at
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for postId, count := range viewCounts {
		_, err = stmt.ExecContext(ctx, postId, count, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetViewCount gets the view count for a single post
func (r *postViewsDBRepository) GetViewCount(ctx context.Context, postId int) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count,
		`SELECT COALESCE(views_count, 0) FROM post_views WHERE post_id = $1`,
		postId,
	)
	if err != nil {
		return 0, nil // Return 0 if not found
	}
	return count, nil
}

// GetMultipleViewCounts gets view counts for multiple posts
func (r *postViewsDBRepository) GetMultipleViewCounts(ctx context.Context, postIds []int) (map[int]int64, error) {
	if len(postIds) == 0 {
		return make(map[int]int64), nil
	}

	query, args, err := sqlx.In(`SELECT post_id, views_count FROM post_views WHERE post_id IN (?)`, postIds)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	viewCounts := make(map[int]int64)
	for rows.Next() {
		var postId int
		var count int64
		if err := rows.Scan(&postId, &count); err != nil {
			return nil, err
		}
		viewCounts[postId] = count
	}

	// Fill in missing post IDs with 0
	for _, postId := range postIds {
		if _, exists := viewCounts[postId]; !exists {
			viewCounts[postId] = 0
		}
	}

	return viewCounts, nil
}
