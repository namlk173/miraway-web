package api

import (
	"go-mirayway/api/middleware"
	"go-mirayway/bootstrap"
	"go-mirayway/handler"
	"go-mirayway/repository"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func InitRoute(app bootstrap.Application) *gin.Engine {
	// Init database and environment variables
	env := *app.Env
	database := app.Mongo.Database(env.DBName)

	// Init Routes
	r := gin.Default()
	r.Use(middleware.AddHeader())

	userRepository := repository.NewUserRepository(database, "user")
	postRepository := repository.NewPostRepository(database, "post")

	userHandler := handler.UserHandler{
		UserRepository: userRepository,
		PostRepository: postRepository,
		Env:            env,
	}

	postHandler := handler.PostHandler{
		PostRepository: postRepository,
		UserRepository: userRepository,
		Env:            env,
	}

	v1 := r.Group("api/v1")
	{
		NewUserApi(v1.Group("user"), userHandler)
		NewPostApi(v1.Group("post"), postHandler)
	}
	r.Use(static.Serve("/upload", static.LocalFile("./upload", true)))

	return r
}
