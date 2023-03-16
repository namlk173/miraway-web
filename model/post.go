package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Content   string             `json:"content" bson:"content" binding:"required""`
	Owner     UserReader         `json:"owner" bson:"owner"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty," bson:"updated_at,omitempty"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
}

type PostRequest struct {
	Title   string `json:"title" bson:"title" binding:"required"`
	Content string `json:"content" bson:"content" binding:"required""`
}

type PostWriter struct {
	Title     string     `json:"title" bson:"title" binding:"required"`
	Content   string     `json:"content" bson:"content" binding:"required""`
	Owner     UserReader `json:"owner" bson:"owner"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

type (
	PostRepository interface {
		Create(ctx context.Context, post *PostWriter) (interface{}, error)
		Delete(ctx context.Context, postID, userID primitive.ObjectID) (int64, error)
		Update(ctx context.Context, postID, userID primitive.ObjectID, post PostRequest) (int64, error)
		Find(ctx context.Context, ID primitive.ObjectID) (*Post, error)
		List(ctx context.Context) ([]Post, error)
		ListPostByUser(ctx context.Context, userID primitive.ObjectID) ([]Post, error)
	}
)
