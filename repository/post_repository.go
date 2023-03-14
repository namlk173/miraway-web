package repository

import (
	"context"
	"go-mirayway/model"
	"go-mirayway/mongodbImplement"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (postRepo *postRepository) Create(ctx context.Context, post *model.PostWriter) error {
	collection := postRepo.db.Collection(postRepo.collection)
	_, err := collection.InsertOne(ctx, post)
	return err
}

func (postRepo *postRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	collection := postRepo.db.Collection(postRepo.collection)
	_, err := collection.DeleteOne(ctx, bson.D{{"_id", ID}})
	return err
}

func (postRepo *postRepository) Update(ctx context.Context, ID primitive.ObjectID, post model.Post) error {
	collection := postRepo.db.Collection(postRepo.collection)
	filter := bson.D{{"_id", ID}}
	updateQuery := bson.D{{"$set", bson.D{{"title", post.Title}, {"content", post.Content}}}}
	_, err := collection.UpdateOne(ctx, filter, updateQuery)
	return err
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
