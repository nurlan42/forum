package server

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func (c *AppContext) Server() {
	// database close
	defer c.db.Close()

	port := flag.String("port", ":8080", "server port")
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)
	mux.HandleFunc("/", c.index)
	mux.HandleFunc("/signup", c.signup)
	mux.HandleFunc("/login", c.login)
	mux.HandleFunc("/logout", c.logout)
	mux.HandleFunc("/newpost", c.showNewPost)
	mux.HandleFunc("/post/", c.showPost)
	mux.HandleFunc("/comment", c.comment)
	mux.HandleFunc("/newcategory", c.newCategory)
	mux.HandleFunc("/postreaction", c.postReaction)
	mux.HandleFunc("/commentreaction", c.commentReaction)

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
