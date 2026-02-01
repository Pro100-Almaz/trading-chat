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

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	uc := &controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Env:         env,
	}

	// USER ROUTES
	group := r.PathPrefix("/user").Subrouter()
	group.HandleFunc("/all", uc.GetUsers).Methods("GET")
	group.HandleFunc("", uc.GetUserById).Methods("GET")
	group.HandleFunc("", uc.UpdateUser).Methods("PUT")
	group.HandleFunc("", uc.DeleteUser).Methods("DELETE")
}
