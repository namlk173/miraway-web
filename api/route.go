package api

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/bootstrap"
	"go-mirayway/handler"
	"go-mirayway/repository"
)

func InitRoute(app bootstrap.Application) *gin.Engine {
	// Init database and environment variables
	env := *app.Env
	database := app.Mongo.Database(env.DBName)

	// Init Routes
	r := gin.Default()

	userRepository := repository.NewUserRepository(database, "user")
	userHandler := handler.UserHandler{
		UserRepository: userRepository,
		Env:            env,
	}

	postRepository := repository.NewPostRepository(database, "post")
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

	return r
}
