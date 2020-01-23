package container

import (
	"context"
	"fmt"
	"log"

	"github.com/stoyaneft/blog/blog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
}

func NewMongoStore() MongoStore {
	return MongoStore{client: nil}
}

func (c *MongoStore) Connect() error {
	var err error
	c.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	return err
}

// GetAll implements blog.Container.
func (c *MongoStore) GetAll() ([]blog.Post, error) {
	if c.client == nil {
		return nil, fmt.Errorf("mongo store is not connected")
	}

	ctx := context.TODO()
	cur, err := c.client.Database("blog").Collection("posts").Find(ctx, bson.D{})
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
	_, err := c.client.Database("blog").Collection("posts").InsertOne(context.TODO(), post)
	return err
}

// Delete implements blog.Container.
func (c *MongoStore) Delete(id int64) error {
	log.Printf("delete mongo recored: %d", id)
	_, err := c.client.Database("blog").Collection("posts").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
