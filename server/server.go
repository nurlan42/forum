package server

import (
	"crypto/tls"
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

	mux.HandleFunc("/category/new", s.auth(s.categoryNew))

	mux.HandleFunc("/post/", s.post)
	mux.HandleFunc("/post/new", s.auth(s.postNew))
	mux.HandleFunc("/post/reaction/", s.auth(s.postReaction))

	mux.HandleFunc("/signin", s.signin)
	mux.HandleFunc("/signup", s.signup)
	mux.HandleFunc("/signout", s.signout)

	mux.HandleFunc("/comment/new", s.auth(s.commentNew))
	mux.HandleFunc("/comment/reaction", s.auth(s.commentReaction))

	mux.HandleFunc("/filter", s.filter)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	cert, err := tls.LoadX509KeyPair("tls/server.crt", "tls/server.key")
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:           *port,
		Handler:        s.limit(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		TLSConfig:      &tls.Config{Certificates: []tls.Certificate{cert}},
		MaxHeaderBytes: 1 << 20,
	}
	s.InfoLog.Printf("Starting server on: https://localhost%v", *port)
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
