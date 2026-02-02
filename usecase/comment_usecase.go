package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"
)

type commentUseCase struct {
	commentRepository repository.CommentRepository
	postRepository    repository.PostRepository
	userRepository    repository.UserRepository
	contextTimeout    time.Duration
}

func NewCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	timeout time.Duration,
) domain.CommentUseCase {
	return &commentUseCase{
		commentRepository: commentRepo,
		postRepository:    postRepo,
		userRepository:    userRepo,
		contextTimeout:    timeout,
	}
}

func (uc *commentUseCase) GetComments(ctx context.Context, postId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	comments, err := uc.commentRepository.GetCommentsByPostId(ctx, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.commentRepository.GetCommentsCount(ctx, postId)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		user, err := uc.userRepository.GetUserById(ctx, comment.UserId)
		if err != nil {
			return nil, err
		}

		responses = append(responses, &domain.CommentResponse{
			Id:        comment.Id,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			Author: domain.Author{
				Id:          user.Id,
				Name:        user.Name,
				AvatarEmoji: user.AvatarEmoji,
			},
		})
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *commentUseCase) CreateComment(ctx context.Context, userId, postId int, request *domain.CreateCommentRequest) (*domain.CommentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if request.Body == "" {
		return nil, errors.New("body is required")
	}

	// Verify post exists
	_, err := uc.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return nil, err
	}

	comment := &domain.Comment{
		UserId: userId,
		PostId: postId,
		Body:   request.Body,
	}

	createdComment, err := uc.commentRepository.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &domain.CommentResponse{
		Id:        createdComment.Id,
		Body:      createdComment.Body,
		CreatedAt: createdComment.CreatedAt,
		Author: domain.Author{
			Id:          user.Id,
			Name:        user.Name,
			AvatarEmoji: user.AvatarEmoji,
		},
	}, nil
}

func (uc *commentUseCase) DeleteComment(ctx context.Context, userId, commentId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	comment, err := uc.commentRepository.GetCommentById(ctx, commentId)
	if err != nil {
		return err
	}

	if comment.UserId != userId {
		return errors.New("you can only delete your own comments")
	}

	return uc.commentRepository.DeleteComment(ctx, commentId)
}
