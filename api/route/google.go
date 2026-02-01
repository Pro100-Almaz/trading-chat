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

func NewGoogleRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	gc := &controller.GoogleController{
		GoogleUseCase: usecase.NewGoogleUseCase(ur, timeout),
		Env:           env,
	}

	r.HandleFunc("/google/login", gc.HandleGoogleLogin).Methods("GET")
	r.HandleFunc("/google/callback", gc.HandleGoogleCallback).Methods("GET")
}
