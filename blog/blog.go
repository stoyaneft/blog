package blog

import "time"

type Post struct {
	ID        int64
	CreatedAt time.Time
	Heading   string
	Author    string
	Content   string
	Likes     int64
	Comments  []Comment
}

type Comment struct {
	Author  string
	Content string
}

type PostContainer interface {
	GetAll() ([]Post, error)
	Insert(*Post) error
	Delete(int64) error
}

type Blog struct {
	posts PostContainer
}

func New(posts PostContainer) *Blog {
	return &Blog{
		posts: posts,
	}
}

func (b *Blog) GetAll() ([]Post, error) {
	return b.posts.GetAll()
}

func (b *Blog) NewPost(post *Post) error {
	return b.posts.Insert(post)
}

func (b *Blog) DeletePost(id int64) error {
	return b.posts.Delete(id)
}
