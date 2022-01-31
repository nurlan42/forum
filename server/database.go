package server

import (
	"database/sql"
	"log"
)

func ConnectDB(driverName string, dbName string) (*AppContext, error) {
	db, err := sql.Open(driverName, dbName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &AppContext{db}, nil
}

func (c *AppContext) CreatePeopleTable() {

	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "people" (
		"user_id"	INTEGER NOT NULL,
		"email"	TEXT NOT NULL UNIQUE,
		"username"	TEXT NOT NULL,
		"password"	BLOB NOT NULL,
		PRIMARY KEY("user_id" AUTOINCREMENT)
	);`)
	CheckErr(err)
	stmt.Exec()
	defer stmt.Close()

}

func (c *AppContext) CreateSessionTable() {
	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "sessions" (
		"user_id"	INTEGER NOT NULL,
		"session_id"	TEXT NOT NULL,
		"last_activity"	DATETIME NOT NULL,
		FOREIGN KEY("user_id") REFERENCES "people"("user_id")
	);`)
	CheckErr(err)
	stmt.Exec()
	defer stmt.Close()
}

func (c *AppContext) CreatePostsTable() {
	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
		"post_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"title"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		"time_creation"	DATETIME NOT NULL,
		PRIMARY KEY("post_id" AUTOINCREMENT),
		FOREIGN KEY("user_id") REFERENCES "people"("user_id")
	);`)
	if err != nil {
		log.Fatalln(err)
	}
	stmt.Exec()
	defer stmt.Close()
}

func (c *AppContext) CreateCommentsTable() {
	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "comments" (
		"comment_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"post_id"	INTEGER NOT NULL,
		"content"	TEXT NOT NULL,
		"time_creation"	DATETIME NOT NULL,
		FOREIGN KEY("user_id") REFERENCES "people"("user_id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		PRIMARY KEY("comment_id")
	);`)
	CheckErr(err)
	stmt.Exec()
	defer stmt.Close()
}

func (c *AppContext) CreateCategoryTable() {
	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "categories" (
		"category_id"	INTEGER NOT NULL UNIQUE,
		"title"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("category_id" AUTOINCREMENT)
	);`)
	CheckErr(err)
	stmt.Exec()
	defer stmt.Close()
}

func (c *AppContext) CreatePostCategory() {
	stmt, err := c.db.Prepare(`CREATE TABLE IF NOT EXISTS "post_category" (
		"pc_id"	INTEGER NOT NULL,
		"post_id"	INTEGER NOT NULL,
		"category_id"	INTEGER NOT NULL,
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		FOREIGN KEY("category_id") REFERENCES "categories"("category_id"),
		PRIMARY KEY("pc_id" AUTOINCREMENT)
	);`)
	CheckErr(err)
	stmt.Exec()
	defer stmt.Close()
}
