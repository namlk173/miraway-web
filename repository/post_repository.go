package repository

import (
	"context"
	"go-mirayway/model"
	"go-mirayway/mongodbImplement"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type postRepository struct {
	db         mongodbImplement.Database
	collection string
}

func NewPostRepository(db mongodbImplement.Database, collection string) model.PostRepository {
	return &postRepository{
		db:         db,
		collection: collection,
	}
}

func (postRepo *postRepository) Create(ctx context.Context, post *model.PostWriter) (interface{}, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	id, err := collection.InsertOne(ctx, post)
	return id, err
}

func (postRepo *postRepository) Delete(ctx context.Context, postID, userID primitive.ObjectID) (int64, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	filter := bson.D{{"_id", postID}, {"owner._id", userID}, {"is_deleted", bson.D{{"$ne", true}}}}
	deleteQuery := bson.D{{"$set", bson.D{{"is_deleted", true}}}}
	res, err := collection.UpdateOne(ctx, filter, deleteQuery)
	return res.MatchedCount, err
}

func (postRepo *postRepository) Update(ctx context.Context, postID, userID primitive.ObjectID, post model.PostRequest) (int64, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	filter := bson.D{{"_id", postID}, {"owner._id", userID}}
	updateQuery := bson.D{{"$set", bson.D{{"title", post.Title}, {"content", post.Content}, {"updated_at", time.Now()}}}}
	res, err := collection.UpdateOne(ctx, filter, updateQuery)
	return res.MatchedCount, err
}

func (postRepo *postRepository) Find(ctx context.Context, ID primitive.ObjectID) (*model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post model.Post
	if err := collection.FindOne(ctx, bson.D{{"_id", ID}}).Decode(&post); err != nil {
		return &model.Post{}, err
	}
	return &post, nil
}

func (postRepo *postRepository) List(ctx context.Context) ([]model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post []model.Post
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return []model.Post{}, err
	}

	if err := cur.All(ctx, &post); err != nil {
		return []model.Post{}, err
	}

	return post, nil
}

func (postRepo *postRepository) ListPostByUser(ctx context.Context, userID primitive.ObjectID) ([]model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post []model.Post
	cur, err := collection.Find(ctx, bson.D{{"author._id", userID}})
	if err != nil {
		return []model.Post{}, err
	}

	if err := cur.All(ctx, &post); err != nil {
		return []model.Post{}, err
	}

	return post, nil
}
