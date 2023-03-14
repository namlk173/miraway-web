package api

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/bootstrap"
	"go-mirayway/handler"
	"go-mirayway/mongodbImplement"
	"go-mirayway/repository"
	"time"
)

func InitRoute(database mongodbImplement.Database, env bootstrap.Env) *gin.Engine {
	r := gin.Default()

	userRepository := repository.NewUserRepository(database, "user")
	userHandler := handler.UserHandler{
		UserRepository: userRepository,
		Env:            env,
	}

	postRepository := repository.NewPostRepository(database, "post")
	postHandler := handler.PostHandler{
		PostRepository: postRepository,
		Timeout:        time.Duration(env.ContextTimeout),
	}

	v1 := r.Group("api/v1")
	{
		NewUserApi(v1.Group("user"), userHandler)
		NewPostApi(v1.Group("post"), postHandler)
	}

	return r
}
