package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *AppContext) post(w http.ResponseWriter, r *http.Request) {

	if !c.alreadyLogIn(r) {
		ErrorPage(w, http.StatusForbidden, "please, log-in first")
		return
	}

	str := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(str)
	CheckErr(err)
	// query to get the post from db
	rows, err := c.db.Query(`SELECT post_id, people.username, title, content, 
	time_creation FROM posts INNER JOIN people on posts.user_id = people.user_id WHERE post_id = ?`, postID)
	CheckErr(err)

	// putting data into struct from db
	var p Post
	for rows.Next() {
		var t time.Time
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &t)
		p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		CheckErr(err)
	}
	// retrieve categories from db
	rs, err := c.db.Query(`SELECT categories.title FROM post_category INNER JOIN categories 
	on post_category.category_id = categories.category_id WHERE post_category.post_id = ?;`, postID)

	for rs.Next() {
		var title string
		err := rs.Scan(&title)
		CheckErr(err)
		p.Categories = append(p.Categories, title)
	}

	// retrive comments from db
	rows, err = c.db.Query(`SELECT people.username, content, 
	time_creation FROM comments INNER JOIN PEOPLE on comments.user_id = people.user_id WHERE post_id = ?`, postID)
	CheckErr(err)

	for rows.Next() {
		var t time.Time
		var comment Comment
		err := rows.Scan(&comment.Author, &comment.Content, &t)
		CheckErr(err)

		comment.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")

		p.Comments = append(p.Comments, comment)
	}

	err = tmpl.ExecuteTemplate(w, "post.html", p)
	if err != nil {
		log.Fatal(err)
	}
}
