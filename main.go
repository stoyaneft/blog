package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/stoyaneft/blog/blog"
	"github.com/stoyaneft/blog/container"
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
		log.Panicln("failed to get posts: %w", err)
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
}

func (s *rest) deletePost(w http.ResponseWriter, r *http.Request) {
	idParams, ok := r.URL.Query()["id"]
	if !ok || len(idParams) == 0 {
		log.Printf("id is required")
		return
	}
	id, err := strconv.ParseInt(idParams[0], 10, 64)
	if err != nil {
		log.Printf("wrong id request: %s", idParams[0])
		return
	}
	err = s.blog.DeletePost(id)
	if err != nil {
		log.Printf("failed to delete post: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Post deleted: %d", id)
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	container := container.NewMongoStore()
	container.Connect()
	blog := blog.New(&container)

	rest := rest{
		server: server,
		mux:    mux,
		blog:   blog,
	}
	rest.Run()
}
