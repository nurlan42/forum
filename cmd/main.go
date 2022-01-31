package main

import (
	"fmt"
	"log"

	"git.01.alem.school/Nurlan/forum.git/server/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := server.ConnectDB("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatalln(err)
	}

	// create all the necessary tables
	db.CreatePeopleTable()
	db.CreateSessionTable()
	db.CreatePostsTable()
	db.CreateCommentsTable()
	db.CreateCategoryTable()
	db.CreatePostCategory()

	fmt.Println("==== database created successfully ====")

	// run server
	db.Server()
}
