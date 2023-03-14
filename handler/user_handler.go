package handler

import (
	"context"
	"fmt"
	"go-mirayway/bootstrap"
	"go-mirayway/model"
	"go-mirayway/util/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)
import "github.com/gin-gonic/gin"

type UserHandler struct {
	UserRepository model.UserRepository
	Env            bootstrap.Env
}

func (userHandler *UserHandler) Signup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	var user model.SignupRequest
	if err := c.ShouldBind(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "some data invalid"})
		return
	}

	if _, err := userHandler.UserRepository.GetUserByEmail(ctx, user.Email); err == nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "email has exist"})
		return
	}

	//if err := util.ValidationPassword(user.Password); err != nil {
	//	c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
	//	return
	//}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	user.Password = string(hashPassword)
	if err := userHandler.UserRepository.CreateUser(ctx, &user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusOK, model.Message{Message: "register successfully."})
	}
}

func (userHandler *UserHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	var loginRequest model.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	userReal, err := userHandler.UserRepository.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: fmt.Sprintf("not found user with email: %v", loginRequest.Email)})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userReal.Password), []byte(loginRequest.Password)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "wrong password"})
		return
	}

	accessToken, err := token.CreateAccessToken(userReal, userHandler.Env.AccessTokenSecret, userHandler.Env.AccessTokenExpiry)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	refreshToken, err := token.CreateRefreshToken(userReal, userHandler.Env.RefreshTokenSecret, userHandler.Env.RefreshTokenExpiry)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	tokenResponse := model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.IndentedJSON(http.StatusAccepted, tokenResponse)
}

func (userHandler *UserHandler) Profile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	idStr, exist := c.GetQuery("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "ID not given"})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "id not true"})
		return
	}

	user, err := userHandler.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (userHandler *UserHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	var password model.Password
	if err := c.ShouldBind(&password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "password invalid"})
		return
	}

	idStr, exist := c.GetQuery("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "not given id"})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	if err := userHandler.UserRepository.UpdatePassword(ctx, id, password.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, model.Message{Message: "change password successful"})
}
