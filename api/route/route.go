package route

import (
	"time"

	"github.com/Pro100-Almaz/trading-chat/api/middleware"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, redisClient *redis.Client, r *mux.Router) {
	public := r.PathPrefix("/api").Subrouter()
	protectedRouter := r.PathPrefix("/api").Subrouter()

	// Initialize token blacklist repository
	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)

	// Middleware to verify AccessToken
	// pass env to middleware
	public.Use(middleware.LoggerMiddleware)
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret, tokenBlacklistRepo))
	protectedRouter.Use(middleware.LoggerMiddleware)

	NewEmojiRouter(public)
	NewGoogleRouter(env, timeout, db, public)
	NewSignupRouter(env, timeout, db, public)
	NewLoginRouter(env, timeout, db, public)
	NewRefreshTokenRouter(env, timeout, db, public)
	NewLogoutRouter(env, timeout, db, protectedRouter)
	NewUserRouter(env, timeout, db, protectedRouter)
	NewVerificationRouter(env, timeout, db, public)
	NewPostRouter(env, timeout, db, redisClient, protectedRouter)
	NewFollowerRouter(env, timeout, db, protectedRouter)
}
