package usecase

import (
	"context"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"
)

type likeUseCase struct {
	likeRepository repository.LikeRepository
	postRepository repository.PostRepository
	contextTimeout time.Duration
}

func NewLikeUseCase(
	likeRepo repository.LikeRepository,
	postRepo repository.PostRepository,
	timeout time.Duration,
) domain.LikeUseCase {
	return &likeUseCase{
		likeRepository: likeRepo,
		postRepository: postRepo,
		contextTimeout: timeout,
	}
}

func (uc *likeUseCase) LikePost(ctx context.Context, userId, postId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Verify post exists
	_, err := uc.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return err
	}

	return uc.likeRepository.LikePost(ctx, userId, postId)
}

func (uc *likeUseCase) UnlikePost(ctx context.Context, userId, postId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.likeRepository.UnlikePost(ctx, userId, postId)
}
