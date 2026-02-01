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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	lc := &controller.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}

	r.HandleFunc("/login", lc.Login).Methods("POST")
}
