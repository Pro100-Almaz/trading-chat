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

func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUseCase: usecase.NewRefreshTokenUseCase(ur, timeout),
		Env:                 env,
	}

	r.HandleFunc("/refresh_token", rtc.RefreshToken).Methods("POST")
}
