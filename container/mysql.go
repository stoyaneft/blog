package container

import (
	"context"
	"fmt"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stoyaneft/blog/blog"
	"go.mongodb.org/mongo-driver/bson"
)

type MySQLStore struct {
	client *sql.DB
}

func NewMySQLStore() MySQLStore {
	return MySQLStore{client: nil}
}

// Connect implements blog.Container.
func (c *MySQLStore) Connect() error {
	var err error
	c.client, err = sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		return err
	}
	defer c.client.Close()
	return nil
}

// GetAll implements blog.Container.
func (c *MySQLStore) GetAll() ([]blog.Post, error) {
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
func (c *MySQLStore) Insert(post *blog.Post) error {
	_, err := c.client.Database("blog").Collection("posts").InsertOne(context.TODO(), post)
	return err
}

// Delete implements blog.Container.
func (c *MySQLStore) Delete(id int64) error {
	log.Printf("delete mongo recored: %d", id)
	_, err := c.client.Database("blog").Collection("posts").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
