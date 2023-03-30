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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	idAny, exits := c.Get("x-user-id")
	if !exits {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", idAny))
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	user, err := postHandler.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	var fileNameSave string

	if post.File != nil {
		extension := filepath.Ext(post.File.Filename)
		fileNameSave := "upload/post/" + uuid.New().String() + extension
		if err := c.SaveUploadedFile(post.File, fileNameSave); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}

	postWriter := model.PostWriter{
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

	idStr := c.Request.URL.Query().Get("_id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Message{Message: "id not true"})
		return
	}

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

	userIdAny, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userIdAny))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "Unauthorized"})
		return
	}

	postIdStr := c.Request.URL.Query().Get("_id")
	postID, err := primitive.ObjectIDFromHex(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: "_id not true"})
		return
	}

	var postRequest model.PostRequest
	if err := c.ShouldBind(&postRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "Required data"})
		return
	}

	if err := util.ValidatePost(postRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	res, err := postHandler.PostRepository.Update(ctx, postID, userID, postRequest)
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

	userIdAny, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userIdAny))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "Unauthorized"})
		return
	}

	postIdStr := c.Request.URL.Query().Get("_id")
	postID, err := primitive.ObjectIDFromHex(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: "_id not true"})
		return
	}

	res, err := postHandler.PostRepository.Delete(ctx, postID, userID)
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

	userIdStr := c.Request.URL.Query().Get("_id")
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Message{Message: "user not found"})
		return
	}

	posts, err := postHandler.PostRepository.ListPostByUser(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, posts)
}
