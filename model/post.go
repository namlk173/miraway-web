package model

import (
	"context"
	"mime/multipart"
	"time"
)

type Post struct {
	ID        string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string     `json:"title" bson:"title" binding:"required"`
	Content   string     `json:"content" bson:"content" binding:"required"`
	ImageURL  string     `json:"image,omitempty" bson:"image,omitempty"`
	Owner     UserReader `json:"owner" bson:"owner"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty," bson:"updated_at,omitempty"`
	IsDeleted bool       `json:"is_deleted" bson:"is_deleted"`
}

type PostRequest struct {
	Title    string                `form:"title" bson:"title" binding:"required"`
	Content  string                `form:"content" bson:"content" binding:"required"`
	ImageURL string                `form:"image" json:"image,omitempty" bson:"image,omitempty"`
	File     *multipart.FileHeader `form:"file"`
}

type (
	PostRepository interface {
		Create(ctx context.Context, post *Post) (interface{}, error)
		Delete(ctx context.Context, postID, userID string) (int64, error)
		Update(ctx context.Context, postID, userID string, post PostRequest) (int64, error)
		UpdateOwner(ctx context.Context, owner *UserReader) error
		Find(ctx context.Context, ID string) (*Post, error)
		List(ctx context.Context, skip, limit int64) ([]Post, error)
		ListPostByUser(ctx context.Context, userID string) ([]Post, error)
	}
)
