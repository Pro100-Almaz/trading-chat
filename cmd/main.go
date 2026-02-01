package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Pro100-Almaz/trading-chat/api/route"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/utils"

	_ "github.com/Pro100-Almaz/trading-chat/docs"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Trading Chat API
// @version 1.0
// @description REST API for Trading Chat application with user authentication and management.

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.Postgres
	defer app.CloseDBConnection()

	utils.MigrateDB(db)

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := mux.NewRouter()

	// Swagger documentation route
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	route.Setup(env, timeout, db, r)

	srv := &http.Server{
		Addr:         env.ServerAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	log.Info("server started")
	log.Info("Swagger UI available at http://localhost:8080/swagger/index.html")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info("shutting down")
	os.Exit(0)
}
