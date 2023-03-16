package middleware

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/model"
	"go-mirayway/util/token"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		str := strings.Split(authHeader, " ")
		if len(str) == 2 {
			authToken := str[1]
			authorized, err := token.IsAuthorized(authToken, secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, model.Message{Message: err.Error()})
				c.Abort()
				return
			}

			if authorized {
				userID, err := token.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, model.Message{Message: err.Error()})
					c.Abort()
					return
				}
				c.Set("x-user-id", userID)
				c.Next()
				return
			}

			c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
			c.Abort()
			return
		}

		c.JSON(http.StatusUnauthorized, model.Message{Message: "Not authorized"})
		c.Abort()
	}
}
