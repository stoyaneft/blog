package container

import (
	"context"
	"fmt"

	"github.com/stoyaneft/blog/blog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	opts   MongoOptions
	client *mongo.Client
}

type MongoOptions struct {
	URI string
}

func NewMongoStore(opts MongoOptions) MongoStore {
	return MongoStore{client: nil, opts: opts}
}

func (c *MongoStore) Init() error {
	var err error
	c.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(c.opts.URI))
	return err
}

// GetAll implements blog.Container.
func (c *MongoStore) GetAll() ([]blog.Post, error) {
	if c.client == nil {
		return nil, fmt.Errorf("mongo store is not connected")
	}

	ctx := context.TODO()
	cur, err := c.collection().Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain posts: %w", err)
	}
	defer cur.Close(ctx)

	posts := []blog.Post{}
	for cur.Next(ctx) {
		var result blog.Post
		err := cur.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to decode post: %w", err)
		}
		posts = append(posts, result)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}
	return posts, nil
}

// Insert implements blog.Container.
func (c *MongoStore) Insert(post *blog.Post) error {
	if c.client == nil {
		return fmt.Errorf("mongo store is not connected")
	}
	_, err := c.collection().InsertOne(context.TODO(), post)
	return err
}

// Delete implements blog.Container.
func (c *MongoStore) Delete(id int64) error {
	if c.client == nil {
		return fmt.Errorf("mongo store is not connected")
	}

	_, err := c.collection().DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (c *MongoStore) collection() *mongo.Collection {
	return c.client.Database(database).Collection(collection)
}

var database string = "blog"
var collection string = "posts"
