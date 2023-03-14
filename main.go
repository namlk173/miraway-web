package main

import (
	"go-mirayway/api"
	"go-mirayway/bootstrap"
	"log"
)

func main() {
	app := bootstrap.NewApp()
	defer app.CloseDatabase()

	database := app.Mongo.Database(app.Env.DBName)
	env := *app.Env

	r := api.InitRoute(database, env)
	log.Fatalln(r.Run(":8080"))

	//env := bootstrap.NewEnv()
	//
	//accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjY0MTA0OWNmMWU5NGU3ODc4NTUyNDU5YyIsInVzZXJuYW1lIjoibmFtbGsiLCJlbWFpbCI6Imxla2hhY25hbS5wZXJzb25hbEBnbWFpbC5jb20iLCJleHAiOjE2Nzg3OTc1OTl9.lPjNOFqbg6TwG5KocZVcNEAoB0dQvmEWC53ZZ4WdrMQ"
	//refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjY0MTA0OWNmMWU5NGU3ODc4NTUyNDU5YyIsImV4cCI6MTY3OTM5NTE5OX0.5yp6LxzXGimvq-EvxS3NrfFa-P3Zvv2sPRD0oOnqI-A"
	//
	//fmt.Println("Authorized")
	//fmt.Println(token.IsAuthorized(accessToken, env.AccessTokenSecret))
	//fmt.Println(token.IsAuthorized(refreshToken, env.RefreshTokenSecret))
	//
	//fmt.Println("ExtractID")
	//fmt.Println(token.ExtractIDFromToken(accessToken, env.AccessTokenSecret))
	//fmt.Println(token.ExtractIDFromToken(refreshToken, env.RefreshTokenSecret))

}
