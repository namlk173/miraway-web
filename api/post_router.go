package api

import (
	"go-mirayway/api/middleware"
	"go-mirayway/handler"

	"github.com/gin-gonic/gin"
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
	s.GET("/all", s.PostHandler.ListAllPost)
	s.GET("/detail", s.PostHandler.GetPostByID)
	s.POST("/create", s.PostHandler.Create)
	s.PUT("/update", s.UpdatePost)
	s.DELETE("/delete", s.DeletePost)
	s.GET("/all/user", s.ListPostByUser)
	return s
}
