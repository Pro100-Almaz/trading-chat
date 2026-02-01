package route

import (
	"time"

	"github.com/Pro100-Almaz/trading-chat/api/middleware"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	public := r.PathPrefix("/api").Subrouter()
	protectedRouter := r.PathPrefix("/api").Subrouter()

	// Middleware to verify AccessToken
	// pass env to middleware
	public.Use(middleware.LoggerMiddleware)
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	protectedRouter.Use(middleware.LoggerMiddleware)

	NewEmojiRouter(public)
	NewGoogleRouter(env, timeout, db, public)
	NewSignupRouter(env, timeout, db, public)
	NewLoginRouter(env, timeout, db, public)
	NewRefreshTokenRouter(env, timeout, db, public)
	NewUserRouter(env, timeout, db, protectedRouter)
	NewVerificationRouter(env, timeout, db, public)
}
