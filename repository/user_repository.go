package repository

import (
	"context"
	"go-mirayway/model"
	"go-mirayway/mongodbImplement"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct {
	db         mongodbImplement.Database
	collection string
}

func NewUserRepository(db mongodbImplement.Database, collectionName string) model.UserRepository {
	return &userRepository{
		db:         db,
		collection: collectionName,
	}
}

func (userRepo *userRepository) CreateUser(ctx context.Context, request *model.SignupRequest) error {
	collection := userRepo.db.Collection(userRepo.collection)
	_, err := collection.InsertOne(ctx, request)
	return err
}

func (userRepo *userRepository) GetAllUser(ctx context.Context) ([]model.User, error) {
	var users []model.User
	collection := userRepo.db.Collection(userRepo.collection)
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return users, err
	}

	if err := cur.All(ctx, &users); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (userRepo *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	collection := userRepo.db.Collection(userRepo.collection)
	var user model.User
	if err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user); err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

func (userRepo *userRepository) GetUserByID(ctx context.Context, ID primitive.ObjectID) (*model.UserReader, error) {
	collection := userRepo.db.Collection(userRepo.collection)
	var user model.UserReader
	if err := collection.FindOne(ctx, bson.D{{"_id", ID}}).Decode(&user); err != nil {
		return &model.UserReader{}, err
	}

	return &user, nil
}

func (userRepo *userRepository) UpdateUser(ctx context.Context, ID primitive.ObjectID, user *model.UserReader) error {
	collection := userRepo.db.Collection(userRepo.collection)
	filter := bson.D{{"_id", ID}}
	updateQuery := bson.D{{"$set", bson.D{{"firstname", user.FirstName}, {"surname", user.SurName}, {"phone", user.MobilePhone}, {"address1", user.Address1}, {"address2", user.Address2}, {"education", user.Education}, {"country", user.Counttry}, {"state", user.State}, {"avatar_url", user.AvatarURL}}}}
	_, err := collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return err
	}

	return nil
}

func (userRepo *userRepository) UpdatePassword(ctx context.Context, ID primitive.ObjectID, password string) error {
	collection := userRepo.db.Collection(userRepo.collection)
	filter := bson.D{{"_id", ID}}
	updateQuery := bson.D{{"$set", bson.D{{"password", password}}}}
	_, err := collection.UpdateOne(ctx, filter, updateQuery)

	return err
}
