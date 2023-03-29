package main

import (
	"go-mirayway/api"
	"go-mirayway/bootstrap"
	"log"
)

func main() {
	app := bootstrap.NewApp()
	defer app.Close()

	// database := app.Mongo.Database(app.Env.DBName)
	// env := *app.Env

	r := api.InitRoute(*app)
	log.Fatalln(r.Run(":8080"))

}
