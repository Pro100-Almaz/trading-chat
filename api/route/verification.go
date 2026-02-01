package route

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/Pro100-Almaz/trading-chat/api/controller"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/internal/email"
	"github.com/Pro100-Almaz/trading-chat/repository"
	"github.com/Pro100-Almaz/trading-chat/usecase"
)

func NewVerificationRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	ur := repository.NewUserRepository(db)
	vr := repository.NewVerificationRepository(db)
	es := email.NewEmailService(env)

	vc := &controller.VerificationController{
		VerificationUseCase: usecase.NewVerificationUseCase(ur, vr, es, timeout),
	}

	r.HandleFunc("/verify-email", vc.VerifyEmail).Methods("POST")
	r.HandleFunc("/resend-verification", vc.ResendVerificationCode).Methods("POST")
}
