package repository

import (
	"context"
	"fmt"
	"go-mirayway/model"
	"go-mirayway/mongodbImplement"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (postRepo *postRepository) Create(ctx context.Context, post *model.Post) (interface{}, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	id, err := collection.InsertOne(ctx, post)
	return id, err
}

func (postRepo *postRepository) Delete(ctx context.Context, postID, userID string) (int64, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	filter := bson.D{{"_id", postID}, {"owner._id", userID}, {"is_deleted", bson.D{{"$ne", true}}}}
	deleteQuery := bson.D{{"$set", bson.D{{"is_deleted", true}}}}
	res, err := collection.UpdateOne(ctx, filter, deleteQuery)
	return res.MatchedCount, err
}

func (postRepo *postRepository) Update(ctx context.Context, postID, userID string, post model.PostRequest) (int64, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	filter := bson.D{{"_id", postID}, {"owner._id", userID}, {"is_deleted", bson.D{{"$ne", true}}}}
	updateQuery := bson.D{{"$set", bson.D{{"title", post.Title}, {"image", post.ImageURL}, {"content", post.Content}, {"updated_at", time.Now()}}}}
	res, err := collection.UpdateOne(ctx, filter, updateQuery)
	return res.MatchedCount, err
}

func (postRepo *postRepository) UpdateOwner(ctx context.Context, owner *model.UserReader) error {
	collection := postRepo.db.Collection(postRepo.collection)
	fmt.Println("avatar_url", owner.AvatarURL)
	filter := bson.D{{"owner._id", owner.ID}, {"is_deleted", bson.D{{"$ne", true}}}}
	updateQuery := bson.D{{"$set", bson.D{{"owner", owner}}}}
	res, err := collection.UpdateMany(ctx, filter, updateQuery)
	fmt.Println("res", res)
	return err
}

func (postRepo *postRepository) Find(ctx context.Context, ID string) (*model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post model.Post
	if err := collection.FindOne(ctx, bson.D{{"_id", ID}}).Decode(&post); err != nil {
		return &model.Post{}, err
	}
	return &post, nil
}

func (postRepo *postRepository) List(ctx context.Context, skip, limit int64) ([]model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post []model.Post
	skipStage := bson.D{{"$skip", skip}}
	limitStage := bson.D{{"$limit", limit}}
	matchStage := bson.D{{"$match", bson.D{{"is_deleted", bson.D{{"$ne", true}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"created_at", -1}}}}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage})
	if err != nil {
		return []model.Post{}, err
	}

	if err := cur.All(ctx, &post); err != nil {
		return []model.Post{}, err
	}

	return post, nil
}

func (postRepo *postRepository) ListPostByUser(ctx context.Context, userID string) ([]model.Post, error) {
	collection := postRepo.db.Collection(postRepo.collection)
	var post []model.Post
	cur, err := collection.Find(ctx, bson.D{{"owner._id", userID}, {"is_deleted", bson.D{{"$ne", true}}}})
	if err != nil {
		return []model.Post{}, err
	}

	if err := cur.All(ctx, &post); err != nil {
		return []model.Post{}, err
	}

	return post, nil
}
