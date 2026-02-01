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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	sc := controller.SignupController{
		SignupUseCase: usecase.NewSignupUseCase(ur, timeout),
		Env:           env,
	}

	r.HandleFunc("/signup", sc.Signup).Methods("POST")
}
