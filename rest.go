package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"blog/blog"
	"blog/container"
)

type rest struct {
	server *http.Server
	mux    *http.ServeMux
	blog   *blog.Blog
}

func (s *rest) Run() error {
	s.mux.HandleFunc("/", s.handleMain)
	s.mux.HandleFunc("/post", s.createPost)
	s.mux.HandleFunc("/delete", s.deletePost)

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start service on port %s: %w", s.server.Addr, err)
	}
	return nil
}

func (s *rest) handleMain(w http.ResponseWriter, r *http.Request) {
	posts, err := s.blog.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Posts: %+v", posts)
}

func (s *rest) createPost(w http.ResponseWriter, r *http.Request) {
	post := blog.Post{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("failed to decode json request body: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.blog.NewPost(&post)
	if err != nil {
		log.Printf("failed to insert post: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "New post added: %+v", post)
}

func (s *rest) deletePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Post deleted")
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	container := container.NewInMemory()
	blog := blog.New(&container)

	rest := rest{
		server: server,
		mux:    mux,
		blog:   blog,
	}
	rest.Run()
}
