package handler

import (
	"context"
	"fmt"
	"go-mirayway/bootstrap"
	"go-mirayway/model"
	"go-mirayway/util"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	PostRepository model.PostRepository
	UserRepository model.UserRepository
	Env            bootstrap.Env
}

func (postHandler *PostHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Duration(postHandler.Env.ContextTimeout)*time.Second)
	defer cancel()

	var post model.PostRequest
	if err := c.Bind(&post); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	if err := util.ValidatePost(post); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	id, exits := c.Get("x-user-id")
	if !exits {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	user, err := postHandler.UserRepository.GetUserByID(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	var fileNameSave string

	if post.File != nil {
		extension := filepath.Ext(post.File.Filename)
		fileNameSave = "upload/post/" + uuid.New().String() + extension
		if err := c.SaveUploadedFile(post.File, fileNameSave); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}

	fmt.Println(fileNameSave)

	postWriter := model.Post{
		ID:        fmt.Sprintf("post_%v", uuid.New().String()),
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  fileNameSave,
		Owner:     *user,
		CreatedAt: time.Now(),
	}

	res, err := postHandler.PostRepository.Create(ctx, &postWriter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"_id": res})
}

func (postHandler *PostHandler) GetPostByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(postHandler.Env.ContextTimeout))
	defer cancel()

	id := c.Request.URL.Query().Get("_id")

	post, err := postHandler.PostRepository.Find(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Message{Message: "not found post for this _id"})
		return
	}

	if post.IsDeleted {
		c.JSON(http.StatusNotFound, model.Message{Message: "This post has been deleted"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(postHandler.Env.ContextTimeout))
	defer cancel()

	userId, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	postId := c.Request.URL.Query().Get("_id")

	var postRequest model.PostRequest
	if err := c.Bind(&postRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "Required data"})
		return
	}

	if err := util.ValidatePost(postRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	var fileNameSave = postRequest.ImageURL

	if postRequest.File != nil {
		extension := filepath.Ext(postRequest.File.Filename)
		fileNameSave = "upload/post/" + uuid.New().String() + extension
		if err := c.SaveUploadedFile(postRequest.File, fileNameSave); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}
	postRequest.ImageURL = fileNameSave

	res, err := postHandler.PostRepository.Update(ctx, postId, fmt.Sprintf("%v", userId), postRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	if res == 0 {
		c.JSON(http.StatusBadRequest, model.Message{Message: "your are not this post owner or this post not existing"})
		return
	}

	c.IndentedJSON(http.StatusAccepted, model.Message{Message: "updated post successfully"})
}

func (postHandler *PostHandler) DeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(postHandler.Env.ContextTimeout))
	defer cancel()

	userId, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	postId := c.Request.URL.Query().Get("_id")
	res, err := postHandler.PostRepository.Delete(ctx, postId, fmt.Sprintf("%v", userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	if res == 0 {
		c.JSON(http.StatusBadRequest, model.Message{Message: "your are not this post owner or this post not existing"})
		return
	}

	c.IndentedJSON(http.StatusAccepted, model.Message{Message: "deleted post successfully"})

}

func (postHandler *PostHandler) ListAllPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(postHandler.Env.ContextTimeout))
	defer cancel()

	s := c.Request.URL.Query().Get("skip")
	l := c.Request.URL.Query().Get("limit")
	skip, err := strconv.Atoi(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	limit, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	posts, err := postHandler.PostRepository.List(ctx, int64(skip), int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (postHandler *PostHandler) ListPostByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(postHandler.Env.ContextTimeout))
	defer cancel()

	userId := c.Request.URL.Query().Get("_id")

	posts, err := postHandler.PostRepository.ListPostByUser(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, posts)
}
