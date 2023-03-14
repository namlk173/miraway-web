package api

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/handler"
)

type UserApi struct {
	*gin.RouterGroup
	handler.UserHandler
}

func NewUserApi(group *gin.RouterGroup, handler handler.UserHandler) *UserApi {
	s := &UserApi{
		RouterGroup: group,
		UserHandler: handler,
	}

	s.POST("/signup", s.UserHandler.Signup)
	s.POST("/login", s.UserHandler.Login)
	// NEED MIDDLEWARE FOR THIS
	s.GET("/profile", s.UserHandler.Profile)
	s.PUT("/profile", s.UserHandler.ChangePassword)

	return s
}
