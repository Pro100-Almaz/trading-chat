package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Env      *Env
	Postgres *sqlx.DB
	Redis    *redis.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Postgres = NewPostgreSQLDatabase(app.Env)
	app.Redis = NewRedisClient(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgreSQLConnection(app.Postgres)
	CloseRedisConnection(app.Redis)
}
