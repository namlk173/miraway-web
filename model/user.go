package model

import (
	"context"
	"mime/multipart"
)

var Users []User

type User struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string `json:"username" bson:"username" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type SignupRequest struct {
	UserName string `json:"username" bson:"username" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UserReader struct {
	ID          string                `form:"_id,omitempty" json:"_id" bson:"_id,omitempty"`
	UserName    string                `form:"username" json:"username" bson:"username" binding:"required"`
	Email       string                `form:"email" json:"email" bson:"email" binding:"required,email"`
	FirstName   string                `form:"firstname" json:"firstname,omitempty" bson:"firstname,omitempty"`
	SurName     string                `form:"surname" json:"surname,omitempty" bson:"surname,omitempty"`
	MobilePhone string                `form:"phone" json:"phone,omitempty" bson:"phone,omitempty"`
	Address1    string                `form:"address1" json:"address1,omitempty" bson:"address1,omitempty"`
	Address2    string                `form:"address2" json:"address2,omitempty" bson:"address2,omitempty"`
	Education   string                `form:"education" json:"education,omitempty" bson:"education,omitempty"`
	Country     string                `form:"country" json:"country,omitempty" bson:"country,omitempty"`
	State       string                `form:"state" json:"state,omitempty" bson:"state,omitempty"`
	AvatarURL   string                `form:"avatar_url" json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	AvatarFile  *multipart.FileHeader `form:"avatar_file" json:"-" bson:"avatar_file,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type (
	UserRepository interface {
		CreateUser(ctx context.Context, request *User) error
		UpdateUser(ctx context.Context, ID string, user *UserReader) error
		GetAllUser(ctx context.Context) ([]User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByID(ctx context.Context, ID string) (*UserReader, error)
		UpdatePassword(ctx context.Context, ID string, password string) error
	}
)
