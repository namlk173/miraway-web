package bootstrap

import (
	"go-mirayway/mongodbImplement"
)

type Application struct {
	Env   *Env
	Mongo mongodbImplement.Client
}

func NewApp() *Application {
	env := NewEnv()
	return &Application{
		Env:   env,
		Mongo: NewMongoClient(env),
	}
}

func (app *Application) Close() {
	CLoseMongoClient(app.Mongo)
}
