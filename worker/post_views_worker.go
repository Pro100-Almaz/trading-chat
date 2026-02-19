package worker

import (
	"context"
	"time"

	"github.com/Pro100-Almaz/trading-chat/repository"

	log "github.com/sirupsen/logrus"
)

type PostViewsWorker struct {
	redisRepo repository.PostViewsRedisRepository
	dbRepo    repository.PostViewsDBRepository
	interval  time.Duration
	stopCh    chan struct{}
}

func NewPostViewsWorker(
	redisRepo repository.PostViewsRedisRepository,
	dbRepo repository.PostViewsDBRepository,
	interval time.Duration,
) *PostViewsWorker {
	return &PostViewsWorker{
		redisRepo: redisRepo,
		dbRepo:    dbRepo,
		interval:  interval,
		stopCh:    make(chan struct{}),
	}
}

// Start begins the worker that syncs Redis views to PostgreSQL
func (w *PostViewsWorker) Start() {
	log.Info("Post views worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.syncViewsToDatabase()
		case <-w.stopCh:
			log.Info("Post views worker stopped")
			return
		}
	}
}

// Stop stops the worker
func (w *PostViewsWorker) Stop() {
	close(w.stopCh)
}

// syncViewsToDatabase syncs view counts from Redis to PostgreSQL
func (w *PostViewsWorker) syncViewsToDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// Get all view counts from Redis
	viewCounts, err := w.redisRepo.GetAllViewCounts(ctx)
	if err != nil {
		log.Errorf("Failed to get view counts from Redis: %v", err)
		return
	}

	if len(viewCounts) == 0 {
		log.Debug("No views to sync")
		return
	}

	log.Infof("Syncing %d post views to database", len(viewCounts))

	// Upsert to database
	err = w.dbRepo.UpsertViewCounts(ctx, viewCounts)
	if err != nil {
		log.Errorf("Failed to upsert view counts to database: %v", err)
		return
	}

	// Reset Redis counters after successful sync
	for postId := range viewCounts {
		err = w.redisRepo.ResetViewCount(ctx, postId)
		if err != nil {
			log.Warnf("Failed to reset view count for post %d: %v", postId, err)
		}
	}

	log.Infof("Successfully synced %d post views to database", len(viewCounts))
}
