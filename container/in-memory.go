package container

import (
	"github.com/stoyaneft/blog/blog"
)

type InMemory struct {
	id    int64
	posts map[int64]blog.Post
}

func NewInMemory() InMemory {
	return InMemory{
		id:    int64(1),
		posts: map[int64]blog.Post{},
	}
}

// GetAll implements blog.Container.
func (c *InMemory) GetAll() ([]blog.Post, error) {
	posts := []blog.Post{}
	for _, post := range c.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// Insert implements blog.Container.
func (c *InMemory) Insert(post *blog.Post) error {
	post.ID = c.id
	c.id++
	c.posts[post.ID] = *post
	return nil
}

// Delete implements blog.Container.
func (c *InMemory) Delete(id int64) error {
	delete(c.posts, id)
	return nil
}
