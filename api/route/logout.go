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

func NewLogoutRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	tbr := repository.NewTokenBlacklistRepository(db)
	lc := &controller.LogoutController{
		LogoutUseCase: usecase.NewLogoutUseCase(tbr, timeout),
	}

	r.HandleFunc("/logout", lc.Logout).Methods("POST")
}
