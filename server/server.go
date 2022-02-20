package server

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func (s *AppContext) Server(p string) {
	port := flag.String("port", p, "server port")
	mux := http.NewServeMux()

	// new structure
	mux.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)
	mux.HandleFunc("/", s.index)

	mux.HandleFunc("/category/new", s.categoryNew)

	mux.HandleFunc("/post/", s.post)
	mux.HandleFunc("/post/new", s.postNew)
	mux.HandleFunc("/post/reaction", s.postReaction)

	mux.HandleFunc("/signin", s.signin)
	mux.HandleFunc("/signup", s.signup)
	mux.HandleFunc("/signout", s.signout)

	mux.HandleFunc("/comment/new", s.commentNew)
	mux.HandleFunc("/comment/reaction", s.commentReaction)

	mux.HandleFunc("/filter", s.filter)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	srv := &http.Server{
		Addr:           *port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Starting server on %v\nlink: http://localhost%v", *port, *port)
	log.Fatal(srv.ListenAndServe())
}
