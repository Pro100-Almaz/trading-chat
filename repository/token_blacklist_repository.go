package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokenBlacklistRepository interface {
	AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpiredTokens(ctx context.Context) error
}

type tokenBlacklistRepository struct {
	db *sqlx.DB
}

func NewTokenBlacklistRepository(db *sqlx.DB) TokenBlacklistRepository {
	return &tokenBlacklistRepository{
		db: db,
	}
}

func (r *tokenBlacklistRepository) AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO token_blacklist (token, expires_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING`,
		token, expiresAt,
	)
	return err
}

func (r *tokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM token_blacklist WHERE token = $1 AND expires_at > NOW()`,
		token,
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tokenBlacklistRepository) CleanupExpiredTokens(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM token_blacklist WHERE expires_at <= NOW()`)
	return err
}
