package main

import (
	"fmt"
	"html/template"
	"log"

	"forum/pkg/models/sqlite3"
	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlite3.ConnectDb("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.SqlDb.Close()
	// create all the necessary tables
	db.CreatePeopleTable()
	db.CreateSessionTable()
	db.CreatePostsTable()
	db.CreateCommentsTable()
	db.CreateCategoryTable()
	db.CreatePostCategory()
	db.CreatePostReaction()
	db.CreateCommentReaction()
	fmt.Println("==== database created successfully ====")

	template := template.Must(template.ParseGlob("ui/html/*.html"))
	appCtx := server.NewAppContext(db, nil, template)
	appCtx.Server()
}
