package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type signupUseCase struct {
	userRepository      repository.UserRepository
	verificationUseCase domain.VerificationUseCase
	contextTimeout      time.Duration
}

func NewSignupUseCase(userRepository repository.UserRepository, verificationUseCase domain.VerificationUseCase, timeout time.Duration) domain.SignupUseCase {
	return &signupUseCase{
		userRepository:      userRepository,
		verificationUseCase: verificationUseCase,
		contextTimeout:      timeout,
	}
}

func (su *signupUseCase) SignUp(ctx context.Context, request domain.SignupRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error(err)
		return err
	}

	request.Password = string(encryptedPassword)

	// Use provided emoji or assign a random one
	avatarEmoji := rand.Intn(len(domain.AvatarEmojis))
	if request.AvatarEmoji != nil && domain.IsValidEmojiIndex(*request.AvatarEmoji) {
		avatarEmoji = *request.AvatarEmoji
	}

	now := time.Now()
	user := &domain.User{
		Name:        request.Name,
		Email:       request.Email,
		Password:    request.Password,
		AvatarEmoji: avatarEmoji,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	user, err = su.userRepository.CreateUser(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	// Send verification code to user's email
	err = su.verificationUseCase.SendVerificationCode(ctx, user.Id, user.Email)
	if err != nil {
		log.Error("Failed to send verification code: ", err)
		return err
	}

	return nil
}
