package api

import (
	"go-mirayway/api/middleware"
	"go-mirayway/handler"

	"github.com/gin-gonic/gin"
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
	s.POST("/refresh", s.UserHandler.RefreshToken)
	// NEED MIDDLEWARE FOR THIS
	s.Use(middleware.JwtAuthMiddleware(handler.Env.AccessTokenSecret))
	s.GET("/profile", s.UserHandler.Profile)
	s.PUT("/profile/change", s.UserHandler.ChangeProfile)
	s.PUT("/profile/change-password", s.UserHandler.ChangePassword)
	return s
}
