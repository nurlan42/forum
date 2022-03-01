package main

import (
	"html/template"
	"log"
	"os"
	"time"

	"forum/pkg/models/sqlite3"
	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	bold := "\033[1m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	reset := "\033[0m"
	InfoLogger := log.New(os.Stdout, bold+colorGreen+"INFO:\t "+reset, log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, bold+colorRed+"ERROR: \t"+reset, log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sqlite3.ConnectDb("sqlite3", "forum.sqlite")
	if err != nil {
		ErrorLogger.Fatalln(err)
	}
	defer db.SQLDb.Close()

	// create database tables
	db.CreatePeopleTable()
	db.CreateSessionTable()
	db.CreatePostsTable()
	db.CreateCommentsTable()
	db.CreateCategoryTable()
	db.CreatePostCategory()
	db.CreatePostReaction()
	db.CreateCommentReaction()
	InfoLogger.Println("database created successfully")
	// delete inactive sessions
	
	// go deleteSessions(db)

	port := ":8080"
	template := template.Must(template.ParseGlob("ui/html/*.html"))
	appCtx := &server.AppContext {
		Sqlite3: db,
		InfoLog: InfoLogger,
		ErrorLog: ErrorLogger,
		Template: template,
	}

	appCtx.Server(port)

}

// deleteSessions removes inactive sessions
func deleteSessions(db *sqlite3.Database) {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			db.DeleteInactiveSession()
		}
	}
}
