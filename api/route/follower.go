package route

import (
	"time"

	"github.com/Pro100-Almaz/trading-chat/api/controller"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/repository"
	"github.com/Pro100-Almaz/trading-chat/usecase"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewFollowerRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	followerRepo := repository.NewFollowerRepository(db)
	userRepo := repository.NewUserRepository(db)

	followerUseCase := usecase.NewFollowerUseCase(followerRepo, userRepo, timeout)

	followerController := &controller.FollowerController{
		FollowerUseCase: followerUseCase,
		Env:             env,
	}

	// Follower routes under /users prefix
	usersGroup := r.PathPrefix("/users").Subrouter()
	usersGroup.HandleFunc("/{id}/follow", followerController.Follow).Methods("POST")
	usersGroup.HandleFunc("/{id}/follow", followerController.Unfollow).Methods("DELETE")
	usersGroup.HandleFunc("/{id}/followers", followerController.GetFollowers).Methods("GET")
	usersGroup.HandleFunc("/{id}/following", followerController.GetFollowing).Methods("GET")
}
