package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/internal/email"
	"github.com/Pro100-Almaz/trading-chat/repository"

	log "github.com/sirupsen/logrus"
)

type verificationUseCase struct {
	userRepository         repository.UserRepository
	verificationRepository repository.VerificationRepository
	emailService           *email.EmailService
	contextTimeout         time.Duration
}

func NewVerificationUseCase(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	emailService *email.EmailService,
	timeout time.Duration,
) domain.VerificationUseCase {
	return &verificationUseCase{
		userRepository:         userRepo,
		verificationRepository: verificationRepo,
		emailService:           emailService,
		contextTimeout:         timeout,
	}
}

// GenerateVerificationCode generates a random 6-digit code
func GenerateVerificationCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	code := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", code)
}

func (vu *verificationUseCase) SendVerificationCode(ctx context.Context, userId int, userEmail string) error {
	// Generate 6-digit code
	code := GenerateVerificationCode()

	// Set expiration to 10 minutes from now
	expiresAt := time.Now().Add(10 * time.Minute)

	// Store the code in database
	err := vu.verificationRepository.CreateVerificationCode(ctx, userId, code, expiresAt)
	if err != nil {
		log.Error("Failed to create verification code: ", err)
		return err
	}

	// Send email
	err = vu.emailService.SendVerificationCode(userEmail, code)
	if err != nil {
		log.Error("Failed to send verification email: ", err)
		return domain.ErrFailedToSendEmail
	}

	return nil
}

func (vu *verificationUseCase) VerifyEmail(ctx context.Context, userEmail, code string) error {
	// Get user by email
	user, err := vu.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		log.Error("User not found: ", err)
		return domain.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return domain.ErrUserAlreadyVerified
	}

	// Check verification code
	_, err = vu.verificationRepository.GetVerificationCode(ctx, user.Id, code)
	if err != nil {
		log.Error("Invalid verification code: ", err)
		return domain.ErrInvalidVerificationCode
	}

	// Mark user as verified
	err = vu.verificationRepository.MarkUserAsVerified(ctx, user.Id)
	if err != nil {
		log.Error("Failed to mark user as verified: ", err)
		return err
	}

	// Delete verification codes for this user
	err = vu.verificationRepository.DeleteVerificationCodes(ctx, user.Id)
	if err != nil {
		log.Error("Failed to delete verification codes: ", err)
	}

	return nil
}

func (vu *verificationUseCase) ResendVerificationCode(ctx context.Context, userEmail string) error {
	// Get user by email
	user, err := vu.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		log.Error("User not found: ", err)
		return domain.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return domain.ErrUserAlreadyVerified
	}

	// Send new verification code
	return vu.SendVerificationCode(ctx, user.Id, user.Email)
}
