package repository

import (
	"context"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/jmoiron/sqlx"
)

type VerificationRepository interface {
	CreateVerificationCode(ctx context.Context, userId int, code string, expiresAt time.Time) error
	GetVerificationCode(ctx context.Context, userId int, code string) (*domain.VerificationCode, error)
	DeleteVerificationCodes(ctx context.Context, userId int) error
	MarkUserAsVerified(ctx context.Context, userId int) error
}

type verificationRepository struct {
	db *sqlx.DB
}

func NewVerificationRepository(db *sqlx.DB) VerificationRepository {
	return &verificationRepository{
		db: db,
	}
}

func (r *verificationRepository) CreateVerificationCode(ctx context.Context, userId int, code string, expiresAt time.Time) error {
	// Delete any existing codes for this user first
	_, err := r.db.ExecContext(ctx, `DELETE FROM verification_codes WHERE user_id = $1`, userId)
	if err != nil {
		return err
	}

	// Insert new code
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO verification_codes (user_id, code, expires_at) VALUES ($1, $2, $3)`,
		userId, code, expiresAt,
	)
	return err
}

func (r *verificationRepository) GetVerificationCode(ctx context.Context, userId int, code string) (*domain.VerificationCode, error) {
	var vc domain.VerificationCode
	err := r.db.GetContext(ctx, &vc,
		`SELECT * FROM verification_codes WHERE user_id = $1 AND code = $2 AND expires_at > NOW()`,
		userId, code,
	)
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *verificationRepository) DeleteVerificationCodes(ctx context.Context, userId int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM verification_codes WHERE user_id = $1`, userId)
	return err
}

func (r *verificationRepository) MarkUserAsVerified(ctx context.Context, userId int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET is_verified = TRUE, updated_at = NOW() WHERE id = $1`,
		userId,
	)
	return err
}
