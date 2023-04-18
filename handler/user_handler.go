package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-mirayway/bootstrap"
	"go-mirayway/model"
	"go-mirayway/util"
	"go-mirayway/util/token"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepository model.UserRepository
	PostRepository model.PostRepository
	Env            bootstrap.Env
}

// Signup Function handler: Create a new account for user
// -> GET information from user
// -> Validate information (email, username, password)
// -> Generate hash password for user -> store user information to the database
func (userHandler *UserHandler) Signup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	var user model.SignupRequest
	if err := c.ShouldBind(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "some data invalid"})
		return
	}

	user.Email = strings.ToLower(user.Email)
	if _, err := userHandler.UserRepository.GetUserByEmail(ctx, user.Email); err == nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "email has exist"})
		return
	}

	if err := util.ValidateUsername(user.UserName); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	if err := util.ValidatePassword(user.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	user.Password = string(hashPassword)
	userCreate := model.User{
		ID:       fmt.Sprintf("user_%v", uuid.New().String()),
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	}
	if err := userHandler.UserRepository.CreateUser(ctx, &userCreate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, model.Message{Message: "register successfully."})

}

// Login Function handler: Login User using Email sand Password
// -> Get information from user( Email and Password)
// -> Find User in database that have Email equal Email of user has entered
// -> Compare password which user entered and hash password of user
// -> Generate access token and refresh token for user
func (userHandler *UserHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	var loginRequest model.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	loginRequest.Email = strings.ToLower(loginRequest.Email)
	userReal, err := userHandler.UserRepository.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: fmt.Sprintf("not found user with email: %v", loginRequest.Email)})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userReal.Password), []byte(loginRequest.Password)); err != nil {
		c.IndentedJSON(http.StatusUnauthorized, model.Message{Message: "wrong password"})
		return
	}

	accessToken, err := token.CreateAccessToken(&model.UserReader{
		ID:       userReal.ID,
		UserName: userReal.UserName,
		Email:    userReal.Email,
	}, userHandler.Env.AccessTokenSecret, userHandler.Env.AccessTokenExpiry)
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

// Profile function handler: Get user self Profile;
// -> Get id From header of request -> id: (any type)
// -> Create ObjectID from id get before
// -> Get user from database by id
func (userHandler *UserHandler) Profile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	id, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		c.Abort()
		return
	}

	user, err := userHandler.UserRepository.GetUserByID(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// ChangePassword function handler: Change password for user
// -> Get id From header of request -> id: (any type)
// -> Create ObjectID from id get before
// -> Get new password from user
// -> validate password which user just entered
// -> Hashing password for user
// -> Update password for user have id equal id that we extracted before
func (userHandler *UserHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*time.Duration(userHandler.Env.ContextTimeout))
	defer cancel()

	id, exist := c.Get("x-user-id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "Unauthorized"})
		return
	}

	var password = struct {
		Password string `json:"password" bson:"password" binding:"required"`
	}{}
	if err := c.ShouldBind(&password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "password is required"})
		return
	}

	if err := util.ValidatePassword(password.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	if err := userHandler.UserRepository.UpdatePassword(ctx, fmt.Sprintf("%v", id), string(hashPassword)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, model.Message{Message: "change password successful"})
}

func (userHandler *UserHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Duration(userHandler.Env.ContextTimeout)*time.Second)
	defer cancel()

	var refresh = struct {
		Refresh string `json:"refresh" binding:"required"`
	}{}

	if err := c.ShouldBind(&refresh); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "refresh token required"})
		return
	}

	idExtract, err := token.ExtractIDFromToken(refresh.Refresh, userHandler.Env.RefreshTokenSecret)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "refresh token not true"})
		return
	}

	user, err := userHandler.UserRepository.GetUserByID(ctx, idExtract)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "not found user"})
		return
	}

	accessToken, err := token.CreateAccessToken(user, userHandler.Env.AccessTokenSecret, userHandler.Env.AccessTokenExpiry)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	loginResponse := model.LoginResponse{
		RefreshToken: refresh.Refresh,
		AccessToken:  accessToken,
	}
	c.IndentedJSON(http.StatusOK, loginResponse)
}

func (userHandler *UserHandler) ChangeProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Duration(userHandler.Env.ContextTimeout)*time.Second)
	defer cancel()

	id, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusUnauthorized, model.Message{Message: "Unauthorized"})
		c.Abort()
		return
	}

	var user model.UserReader
	if err := c.Bind(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	var userAvatarURL = user.AvatarURL
	if user.AvatarFile != nil {
		extension := filepath.Ext(user.AvatarFile.Filename)
		userAvatarURL = "upload/user/" + uuid.New().String() + extension
		if err := c.SaveUploadedFile(user.AvatarFile, userAvatarURL); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}

	user.AvatarURL = userAvatarURL
	user.AvatarFile = nil

	if err := userHandler.UserRepository.UpdateUser(ctx, fmt.Sprintf("%v", id), &user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	if err := userHandler.PostRepository.UpdateOwner(ctx, &user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, model.Message{Message: "update profile successfully"})
}
