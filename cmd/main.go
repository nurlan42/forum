package main

import (
	"html/template"
	"log"
	"os"

	"forum/pkg/models/sqlite3"
	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	bold := "\033[1m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	reset := "\033[0m"
	InfoLogger := log.New(os.Stdout, bold+colorGreen+"INFO: "+reset, log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, bold+colorRed+"ERROR: "+reset, log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sqlite3.ConnectDb("sqlite3", "forum.db")
	if err != nil {
		ErrorLogger.Fatalln(err)
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
	InfoLogger.Println("==== database created successfully ====")

	// delete inactive sessions
	// ticker := time.NewTicker(5 * time.Second)
	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case <-ticker.C:
	// 			db.DeleteInactiveSession()
	// 		}
	// 	}
	// }()
	port := ":8080"
	template := template.Must(template.ParseGlob("ui/html/*.html"))
	appCtx := server.NewAppContext(db, InfoLogger, ErrorLogger, template)
	appCtx.Server(port)

}
