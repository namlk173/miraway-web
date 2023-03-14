package api

import (
	"github.com/gin-gonic/gin"
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

	s.POST("/create", s.PostHandler.Create)
	return s
}
