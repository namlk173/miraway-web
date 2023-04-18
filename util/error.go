package util

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/model"
)

func AssertNil(c *gin.Context, code int, err error, message ...string) {
	if err != nil {
		if message != nil {
			c.IndentedJSON(code, model.Message{Message: message[0]})
			c.Abort()
			return
		} else {
			c.IndentedJSON(code, model.Message{Message: err.Error()})
			c.Abort()
			return
		}
	}
	return
}
