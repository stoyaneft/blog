package container

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stoyaneft/blog/blog"
)

type MySQLStore struct {
	opts   MySQLOptions
	client *sql.DB
}

type MySQLOptions struct {
	URI string
}

func NewMySQLStore(opts MySQLOptions) MySQLStore {
	return MySQLStore{client: nil, opts: opts}
}

// Connect implements blog.Container.
func (c *MySQLStore) Connect() error {
	var err error
	c.client, err = sql.Open("mysql", c.opts.URI)
	return err
}

// GetAll implements blog.Container.
func (c *MySQLStore) GetAll() ([]blog.Post, error) {
	if c.client == nil {
		return nil, fmt.Errorf("mysql store is not connected")
	}

	posts := []blog.Post{}
	rows, err := c.client.Query("select id, author, content, likes from posts")
	if err != nil {
		return nil, fmt.Errorf("failed to obtains posts from mysql: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result blog.Post
		err := rows.Scan(&result.ID, &result.Author, &result.Content, &result.Likes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}
	return posts, nil
}

// Insert implements blog.Container.
func (c *MySQLStore) Insert(post *blog.Post) error {
	if c.client == nil {
		return fmt.Errorf("mysql store is not connected")
	}

	_, err := c.client.Exec("insert into posts(author, content, likes) VALUES (?, ?, ?)", post.Author, post.Content, post.Likes)
	return err
}

// Delete implements blog.Container.
func (c *MySQLStore) Delete(id int64) error {
	if c.client == nil {
		return fmt.Errorf("mysql store is not connected")
	}

	log.Printf("delete mysql record: %d", id)
	_, err := c.client.Exec("delete from posts where id=?", id)
	return err
}
