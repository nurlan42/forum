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
	db.CreatePostReaction()
	db.CreateCommentReaction()
	fmt.Println("==== database created successfully ====")

	// delete inactive sessions
	// ticker := time.NewTicker(5 * time.Second)
	// // done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		// case <-done:
	// 		// 	return
	// 		case <-ticker.C:
	// 			db.DeleteInactiveSession()
	// 		}
	// 	}
	// }()
	db.Server()
	
}
