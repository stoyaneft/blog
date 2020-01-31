package container

import (
	"fmt"
	"strconv"

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

func (c *MySQLStore) Init() error {
	var err error
	c.client, err = sql.Open("mysql", c.opts.URI)
	return err
}

// GetAll implements blog.Container.
func (c *MySQLStore) GetAll() ([]blog.Post, error) {
	if c.client == nil {
		return nil, fmt.Errorf("mysql store is not initialized")
	}

	posts := []blog.Post{}
	var id int64
	rows, err := c.client.Query("select id, heading, created_at, author, content, likes from posts")
	if err != nil {
		return nil, fmt.Errorf("failed to obtains posts from mysql: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result blog.Post
		err := rows.Scan(&id, &result.Heading, &result.CreatedAt, &result.Author, &result.Content, &result.Likes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		result.ID = strconv.FormatInt(id, 10)
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
		return fmt.Errorf("mysql store is not initialized")
	}

	_, err := c.client.Exec("insert into posts(heading, author, content, likes, created_at) VALUES (?, ?, ?, ?, ?)",
		post.Heading, post.Author, post.Content, post.Likes, post.CreatedAt)
	return err
}

// Delete implements blog.Container.
func (c *MySQLStore) Delete(id string) error {
	if c.client == nil {
		return fmt.Errorf("mysql store is not initialized")
	}

	_, err := c.client.Exec("delete from posts where id=?", id)
	return err
}
