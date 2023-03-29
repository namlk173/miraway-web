package mongodbImplement

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Collection interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	//InsertMany(context.Context, []interface{}) ([]interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	//CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	//Aggregate(context.Context, interface{}) (Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongoImplement.UpdateResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error)
}

type Database interface {
	Collection(string, ...*options.CollectionOptions) Collection
	Client() Client
}

type Client interface {
	Database(string) Database
	Connect(context.Context) error
	Disconnect(context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(context.Context) error
}

type SingleResult interface {
	Decode(interface{}) error
}

type Cursor interface {
	//Close(context.Context) error
	//Next(context.Context) bool
	//Decode(interface{}) error
	All(context.Context, interface{}) error
}

type mongoClient struct {
	client *mongo.Client
}

type mongoDatabase struct {
	database *mongo.Database
}

type mongoCollection struct {
	collection *mongo.Collection
}

type mongoCursor struct {
	cursor *mongo.Cursor
}

type mongoSingleResult struct {
	singleResult *mongo.SingleResult
}

// Initialize mongoCLient and implement Client Interface
func NewClient(url string) (Client, error) {
	clientOptions := options.Client().SetTimeout(time.Second * 2)
	client, err := mongo.NewClient(clientOptions.ApplyURI(url))
	return &mongoClient{client: client}, err
}

// Implement Client interface for mongoClient struct
func (c *mongoClient) Database(DBName string) Database {
	database := c.client.Database(DBName)
	return &mongoDatabase{database: database}
}

func (c *mongoClient) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

func (c *mongoClient) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *mongoClient) StartSession() (mongo.Session, error) {
	return c.client.StartSession()
}

func (c *mongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return c.client.UseSession(ctx, fn)
}

func (c *mongoClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

// Implement Database interface for mongoDatabase struct
func (d *mongoDatabase) Collection(collectionName string, options ...*options.CollectionOptions) Collection {
	return &mongoCollection{collection: d.database.Collection(collectionName, options...)}
}

func (d *mongoDatabase) Client() Client {
	return &mongoClient{client: d.database.Client()}
}

// Implement Collection interface for mongoCollection
func (coll *mongoCollection) Find(ctx context.Context, filter interface{}, options ...*options.FindOptions) (Cursor, error) {
	res, err := coll.collection.Find(ctx, filter, options...)
	return &mongoCursor{cursor: res}, err

}

func (coll *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	return &mongoSingleResult{singleResult: coll.collection.FindOne(ctx, filter)}
}

func (coll *mongoCollection) InsertOne(ctx context.Context, data interface{}) (interface{}, error) {
	res, err := coll.collection.InsertOne(ctx, data)
	return res.InsertedID, err
}

func (coll *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	res, err := coll.collection.DeleteOne(ctx, filter)
	return res.DeletedCount, err
}
func (coll *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return coll.collection.UpdateOne(ctx, filter, update, options...)
}

func (coll *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error) {
	res, err := coll.collection.Aggregate(ctx, pipeline, opts...)
	return &mongoCursor{cursor: res}, err
}

// Implement Cursor interface for mongoCursor struct
func (cur *mongoCursor) All(ctx context.Context, results interface{}) error {
	return cur.cursor.All(ctx, results)
}

// Implement SingleRsult interface for mongoSingleResult struct
func (s *mongoSingleResult) Decode(result interface{}) error {
	return s.singleResult.Decode(result)
}
