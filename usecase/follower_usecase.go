package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"
)

type followerUseCase struct {
	followerRepository repository.FollowerRepository
	userRepository     repository.UserRepository
	contextTimeout     time.Duration
}

func NewFollowerUseCase(
	followerRepo repository.FollowerRepository,
	userRepo repository.UserRepository,
	timeout time.Duration,
) domain.FollowerUseCase {
	return &followerUseCase{
		followerRepository: followerRepo,
		userRepository:     userRepo,
		contextTimeout:     timeout,
	}
}

func (uc *followerUseCase) Follow(ctx context.Context, followerId, followingId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if followerId == followingId {
		return errors.New("you cannot follow yourself")
	}

	// Verify user to follow exists
	_, err := uc.userRepository.GetUserById(ctx, followingId)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.followerRepository.Follow(ctx, followerId, followingId)
}

func (uc *followerUseCase) Unfollow(ctx context.Context, followerId, followingId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.followerRepository.Unfollow(ctx, followerId, followingId)
}

func (uc *followerUseCase) GetFollowers(ctx context.Context, userId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	users, err := uc.followerRepository.GetFollowers(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.followerRepository.GetFollowersCount(ctx, userId)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.FollowUserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, &domain.FollowUserResponse{
			Id:          user.Id,
			Name:        user.Name,
			AvatarEmoji: user.AvatarEmoji,
		})
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *followerUseCase) GetFollowing(ctx context.Context, userId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	users, err := uc.followerRepository.GetFollowing(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.followerRepository.GetFollowingCount(ctx, userId)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.FollowUserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, &domain.FollowUserResponse{
			Id:          user.Id,
			Name:        user.Name,
			AvatarEmoji: user.AvatarEmoji,
		})
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *followerUseCase) IsFollowing(ctx context.Context, followerId, followingId int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.followerRepository.IsFollowing(ctx, followerId, followingId)
}
