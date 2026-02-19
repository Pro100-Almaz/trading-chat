package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type PostViewsRedisRepository interface {
	IncrementViews(ctx context.Context, postIds []int) error
	GetViewCount(ctx context.Context, postId int) (int64, error)
	GetAllViewCounts(ctx context.Context) (map[int]int64, error)
	ResetViewCount(ctx context.Context, postId int) error
}

type postViewsRedisRepository struct {
	redis *redis.Client
}

func NewPostViewsRedisRepository(redis *redis.Client) PostViewsRedisRepository {
	return &postViewsRedisRepository{
		redis: redis,
	}
}

// IncrementViews increments view count for multiple posts in Redis
func (r *postViewsRedisRepository) IncrementViews(ctx context.Context, postIds []int) error {
	if len(postIds) == 0 {
		return nil
	}

	pipe := r.redis.Pipeline()
	for _, postId := range postIds {
		key := fmt.Sprintf("post:views:%d", postId)
		pipe.Incr(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetViewCount gets the view count for a single post from Redis
func (r *postViewsRedisRepository) GetViewCount(ctx context.Context, postId int) (int64, error) {
	key := fmt.Sprintf("post:views:%d", postId)
	count, err := r.redis.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// GetAllViewCounts gets all view counts from Redis (for worker to sync to DB)
func (r *postViewsRedisRepository) GetAllViewCounts(ctx context.Context) (map[int]int64, error) {
	// Get all keys matching post:views:*
	keys, err := r.redis.Keys(ctx, "post:views:*").Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return make(map[int]int64), nil
	}

	// Get all values using MGET
	values, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Build result map
	viewCounts := make(map[int]int64)
	for i, key := range keys {
		// Extract post ID from key: "post:views:123" -> 123
		parts := strings.Split(key, ":")
		if len(parts) != 3 {
			continue
		}

		postId, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		if values[i] == nil {
			viewCounts[postId] = 0
		} else if strVal, ok := values[i].(string); ok {
			count, err := strconv.ParseInt(strVal, 10, 64)
			if err == nil {
				viewCounts[postId] = count
			}
		}
	}

	return viewCounts, nil
}

// ResetViewCount resets the view count for a post in Redis (after syncing to DB)
func (r *postViewsRedisRepository) ResetViewCount(ctx context.Context, postId int) error {
	key := fmt.Sprintf("post:views:%d", postId)
	return r.redis.Del(ctx, key).Err()
}
