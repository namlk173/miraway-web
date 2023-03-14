package handler

import (
	"github.com/gin-gonic/gin"
	"go-mirayway/model"
	"net/http"
	"time"
)

type PostHandler struct {
	model.PostRepository
	Timeout time.Duration
}

func (postHandler *PostHandler) Create(c *gin.Context) {
	var post model.PostRequest
	if err := c.ShouldBind(&post); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "data not valid"})
		return
	}

}
