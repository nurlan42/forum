package server

import (
	"forum/pkg/models/sqlite3"
	"html/template"
	"log"
)

type AppContext struct {
	Sqlite3  *sqlite3.Database
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Template *template.Template
}

 
