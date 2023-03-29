package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Users []User

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username" binding:"required"`
	Email    string             `json:"email" bson:"email" binding:"required,email"`
	Password string             `json:"password" bson:"password" binding:"required"`
}

type SignupRequest struct {
	UserName string `json:"username" bson:"username" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UserReader struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName    string             `json:"username" bson:"username" binding:"required"`
	Email       string             `json:"email" bson:"email" binding:"required,email"`
	FirstName   string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	SurName     string             `json:"surname,omitempty" bson:"surname,omitempty"`
	MobilePhone string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Address1    string             `json:"address1,omitempty" bson:"address1,omitempty"`
	Address2    string             `json:"address2,omitempty" bson:"address2,omitempty"`
	Education   string             `json:"education,omitempty" bson:"education,omitempty"`
	Counttry    string             `json:"country,omitempty" bson:"country,omitempty"`
	State       string             `json:"state,omitempty" bson:"state,omitempty"`
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
		CreateUser(ctx context.Context, request *SignupRequest) error
		UpdateUser(ctx context.Context, ID primitive.ObjectID, user *UserReader) error
		GetAllUser(ctx context.Context) ([]User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByID(ctx context.Context, ID primitive.ObjectID) (*UserReader, error)
		UpdatePassword(ctx context.Context, ID primitive.ObjectID, password string) error
	}
)
