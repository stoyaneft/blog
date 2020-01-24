package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	s.mux.HandleFunc("/create", s.handleCreate)
	s.mux.HandleFunc("/post", s.createPost)
	s.mux.HandleFunc("/delete", s.deletePost)

	log.Printf("server is listening at %s\n", s.server.Addr)

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

	t := template.Must(template.ParseFiles("./templates/index.tmpl.html"))
	// fmt.Printf("posts: %+v\n", posts)
	t.Execute(w, posts)
}

func (s *rest) handleCreate(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/create.tmpl.html"))
	// fmt.Printf("posts: %+v\n", posts)
	t.Execute(w, struct{}{})
}

func (s *rest) createPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post := blog.Post{
		Content: r.Form.Get("content"),
		Author:  r.Form.Get("author"),
	}
	defer http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	err := s.blog.NewPost(&post)
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

	container := container.NewMySQLStore(container.MySQLOptions{
		URI: fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/blog", os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
	})
	// container := container.NewMongoStore(container.MongoOptions{
	// 	URI: "mongodb://localhost:27017",
	// })
	err := container.Connect()
	if err != nil {
		log.Fatal("failed to connect to store: %w", err)
	}
	blog := blog.New(&container)

	rest := rest{
		server: server,
		mux:    mux,
		blog:   blog,
	}
	rest.Run()
}
