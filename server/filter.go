package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func (c *AppContext) filter(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/filter" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var allPosts *[]Post
	categories, err := c.readCategories()
	CheckErr(err)

	if r.FormValue("owner") == "yes" {
		allPosts = c.filterByOwner(r)
	} else if r.FormValue("category") != "" {
		categoryID, err := strconv.Atoi(r.FormValue("category"))
		CheckErr(err)
		allPosts = c.filterByCategory(r, categoryID)
	} else if r.FormValue("reaction") != "" {
		reaction, err := strconv.Atoi(r.FormValue("reaction"))
		CheckErr(err)
		allPosts = c.filterByReaction(reaction)
	}

	data := struct {
		AllPosts   *[]Post
		Categories map[string]int
	}{allPosts, categories}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AppContext) filterByOwner(r *http.Request) *[]Post {
	// created by you
	cookie, err := r.Cookie("session")
	CheckErr(err)

	sID, err := c.getSession(cookie.Value)
	CheckErr(err)

	userID := sID[cookie.Value]

	rows, err := c.db.Query(`SELECT posts.post_id, people.username,title, content, 
		time_creation FROM posts INNER JOIN people on posts.user_id = people.user_id WHERE posts.user_id = ?;`, userID)
	CheckErr(err)

	var ps []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.TimeCreation)
		CheckErr(err)
		ps = append(ps, p)
	}

	return &ps
}

func (c *AppContext) filterByCategory(r *http.Request, categoryID int) *[]Post {

	rows, err := c.db.Query(`SELECT posts.post_id, people.username, posts.title,
		posts.content, posts.time_creation FROM posts INNER JOIN people ON posts.user_id =
		people.user_id INNER JOIN post_category ON post_category.post_id = posts.post_id 
		WHERE post_category.category_id = ?`, categoryID)

	// rows, err := c.db.Query(`SELECT posts.post_id, posts.title, posts.content, posts.time_creation FROM
	// 	posts INNER JOIN post_category ON posts.post_id = post_category.post_id WHERE post_category.category_id = ?;`, categoryID)
	CheckErr(err)

	var ps []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.TimeCreation)
		CheckErr(err)
		ps = append(ps, p)
	}

	return &ps
}

func (c *AppContext) filterByReaction(emotion int) *[]Post {
	rows, err := c.db.Query(`SELECT posts.post_id, people.username, posts.title, posts.content, posts.time_creation FROM posts
		INNER JOIN people ON people.user_id = posts.user_id INNER JOIN post_reaction ON post_reaction.post_id = posts.post_id
		WHERE post_reaction.reaction = ?;`, emotion)
	CheckErr(err)

	var ps []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.TimeCreation)
		CheckErr(err)
		ps = append(ps, p)
	}
	fmt.Println(ps)
	return &ps
}
