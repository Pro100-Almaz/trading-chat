package bootstrap

import (
	"github.com/jmoiron/sqlx"
)

type Application struct {
	Env      *Env
	Postgres *sqlx.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Postgres = NewPostgreSQLDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgreSQLConnection(app.Postgres)
}
