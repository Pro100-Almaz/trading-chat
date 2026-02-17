package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/repository"
)

type postUseCase struct {
	postRepository    repository.PostRepository
	userRepository    repository.UserRepository
	likeRepository    repository.LikeRepository
	commentRepository repository.CommentRepository
	contextTimeout    time.Duration
}

func NewPostUseCase(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	timeout time.Duration,
) domain.PostUseCase {
	return &postUseCase{
		postRepository:    postRepo,
		userRepository:    userRepo,
		likeRepository:    likeRepo,
		commentRepository: commentRepo,
		contextTimeout:    timeout,
	}
}

func (uc *postUseCase) GetGlobalFeed(ctx context.Context, userId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	posts, err := uc.postRepository.GetPosts(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.postRepository.GetPostsCount(ctx)
	if err != nil {
		return nil, err
	}

	responses, err := uc.enrichPosts(ctx, posts, userId)
	if err != nil {
		return nil, err
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *postUseCase) GetFollowingFeed(ctx context.Context, userId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	posts, err := uc.postRepository.GetFollowingPosts(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.postRepository.GetFollowingPostsCount(ctx, userId)
	if err != nil {
		return nil, err
	}

	responses, err := uc.enrichPosts(ctx, posts, userId)
	if err != nil {
		return nil, err
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *postUseCase) GetUserPosts(ctx context.Context, currentUserId, targetUserId, limit, offset int) (*domain.PaginatedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	posts, err := uc.postRepository.GetPostsByUserId(ctx, targetUserId, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.postRepository.GetUserPostsCount(ctx, targetUserId)
	if err != nil {
		return nil, err
	}

	responses, err := uc.enrichPosts(ctx, posts, currentUserId)
	if err != nil {
		return nil, err
	}

	return domain.NewPaginatedResponse(responses, total, limit, offset), nil
}

func (uc *postUseCase) GetPostById(ctx context.Context, userId, postId int) (*domain.PostResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	post, err := uc.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return nil, err
	}

	return uc.enrichPost(ctx, post, userId)
}

func (uc *postUseCase) CreatePost(ctx context.Context, userId int, request *domain.CreatePostRequest) (*domain.PostResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if request.Ticker == "" {
		return nil, errors.New("ticker is required")
	}
	if request.Body == "" {
		return nil, errors.New("body is required")
	}

	post := &domain.Post{
		UserId: userId,
		Ticker: request.Ticker,
		Body:   request.Body,
	}

	createdPost, err := uc.postRepository.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return uc.enrichPost(ctx, createdPost, userId)
}

func (uc *postUseCase) DeletePost(ctx context.Context, userId, postId int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	post, err := uc.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return err
	}

	if post.UserId != userId {
		return errors.New("you can only delete your own posts")
	}

	return uc.postRepository.DeletePost(ctx, postId)
}

func (uc *postUseCase) enrichPosts(ctx context.Context, posts []*domain.Post, userId int) ([]*domain.PostResponse, error) {
	responses := make([]*domain.PostResponse, 0, len(posts))
	for _, post := range posts {
		response, err := uc.enrichPost(ctx, post, userId)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (uc *postUseCase) enrichPost(ctx context.Context, post *domain.Post, userId int) (*domain.PostResponse, error) {
	user, err := uc.userRepository.GetUserById(ctx, post.UserId)
	if err != nil {
		return nil, err
	}

	likesCount, err := uc.likeRepository.GetLikesCount(ctx, post.Id)
	if err != nil {
		return nil, err
	}

	commentsCount, err := uc.commentRepository.GetCommentsCount(ctx, post.Id)
	if err != nil {
		return nil, err
	}

	isLiked, err := uc.likeRepository.IsLiked(ctx, userId, post.Id)
	if err != nil {
		return nil, err
	}

	return &domain.PostResponse{
		Id:            post.Id,
		Ticker:        post.Ticker,
		Body:          post.Body,
		CreatedAt:     post.CreatedAt,
		Author:        domain.Author{Id: user.Id, Name: user.Name, AvatarEmoji: user.AvatarEmoji},
		LikesCount:    likesCount,
		CommentsCount: commentsCount,
		IsLiked:       isLiked,
	}, nil
}
