package container

import "blog/blog"

type InMemory struct {
	posts map[int64]*blog.Post
}

func NewInMemory() InMemory {
	return InMemory{
		posts: map[int64]*blog.Post{},
	}
}

// GetAll implements blog.Container.
func (c *InMemory) GetAll() ([]blog.Post, error) {
	posts := make([]blog.Post, len(c.posts))
	for _, post := range posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// Insert implements blog.Container.
func (c *InMemory) Insert(post *blog.Post) error {
	c.posts[post.ID] = post
	return nil
}

// Delete implements blog.Container.
func (c *InMemory) Delete(id int64) error {
	delete(c.posts, id)
	return nil
}
