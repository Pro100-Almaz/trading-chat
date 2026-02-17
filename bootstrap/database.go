package bootstrap

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewPostgreSQLDatabase(env *Env) *sqlx.DB {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	// Debug logging
	log.Infof("Connecting to database: host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		log.Errorf("Connection details: host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)
		log.Fatal(err)
	}

	log.Info("Successfully connected to database")
	return db
}

func ClosePostgreSQLConnection(client *sqlx.DB) {
	if client == nil {
		return
	}

	err := client.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connection to PostgreSQL closed.")
}
