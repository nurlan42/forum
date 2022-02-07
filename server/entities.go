package server

import "database/sql"

type AppContext struct {
	db *sql.DB
}

type User struct {
	ID              int
	Email, UserName string
	Pass            []byte
}

type Post struct {
	ID                     int
	UserID                 int
	Title, Content, Author string
	TimeCreation           string
	Comments               []Comment
	Categories             []string
	CommentNbr             int
}

type Comment struct {
	CommID       int
	Author       string
	TimeCreation string
	Content      string
}
