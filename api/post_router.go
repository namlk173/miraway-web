package api

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/api/middleware"
	"go-mirayway/handler"
)

type PostApi struct {
	*gin.RouterGroup
	handler.PostHandler
}

func NewPostApi(group *gin.RouterGroup, handler handler.PostHandler) *PostApi {
	s := &PostApi{
		RouterGroup: group,
		PostHandler: handler,
	}

	s.Use(middleware.JwtAuthMiddleware(handler.Env.AccessTokenSecret))
	s.GET("/", s.PostHandler.GetPostByID)
	s.POST("/create", s.PostHandler.Create)
	s.PUT("/update", s.UpdatePost)
	s.DELETE("/delete", s.DeletePost)
	return s
}
