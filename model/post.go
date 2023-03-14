package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title" binding:"required"`
	Content string             `json:"content" bson:"content" binding:"required""`
	Author  UserReader         `json:"author" bson:"author"`
}

type PostRequest struct {
	Title   string `json:"title" bson:"title" binding:"required"`
	Content string `json:"content" bson:"content" binding:"required""`
}

type PostWriter struct {
	Title   string     `json:"title" bson:"title" binding:"required"`
	Content string     `json:"content" bson:"content" binding:"required""`
	Author  UserReader `json:"author" bson:"author"`
}

type (
	PostRepository interface {
		Create(ctx context.Context, post *PostWriter) error
		Delete(ctx context.Context, ID primitive.ObjectID) error
		Update(ctx context.Context, ID primitive.ObjectID, post Post) error
		Find(ctx context.Context, ID primitive.ObjectID) (*Post, error)
		List(ctx context.Context) ([]Post, error)
		ListPostByUser(ctx context.Context, userID primitive.ObjectID) ([]Post, error)
	}
)
