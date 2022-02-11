package server

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func (s *AppContext) Server() {
	port := flag.String("port", ":8080", "server port")
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/signup", s.signup)
	mux.HandleFunc("/login", s.login)
	mux.HandleFunc("/logout", s.logout)
	mux.HandleFunc("/newpost", s.showNewPost)
	mux.HandleFunc("/post/", s.showPost)
	mux.HandleFunc("/comment", s.comment)
	mux.HandleFunc("/newcategory", s.newCategory)
	mux.HandleFunc("/postreaction", s.postReaction)
	mux.HandleFunc("/commentreaction", s.commentReaction)
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
