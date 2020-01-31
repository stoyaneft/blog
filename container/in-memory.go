package container

import (
	"github.com/stoyaneft/blog/blog"
	"sync"
)

type InMemory struct {
	id    int64
	posts map[int64]blog.Post
	mutex sync.RWMutex
}

func NewInMemory() InMemory {
	return InMemory{
		id:    int64(1),
		posts: map[int64]blog.Post{},
		mutex: sync.RWMutex{},
	}
}

// GetAll implements blog.Container.
func (c *InMemory) GetAll() ([]blog.Post, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	posts := []blog.Post{}
	for _, post := range c.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// Insert implements blog.Container.
func (c *InMemory) Insert(post *blog.Post) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	post.ID = c.id
	c.id++
	c.posts[post.ID] = *post
	return nil
}

// Delete implements blog.Container.
func (c *InMemory) Delete(id int64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.posts, id)
	return nil
}
