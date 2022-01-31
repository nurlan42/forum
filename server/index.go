package server

import (
	"net/http"
	"text/template"
	"time"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("ui/html/*.html"))

}

func (c *AppContext) index(w http.ResponseWriter, r *http.Request) {

	var allPosts []Post

	rows, err := c.db.Query(`SELECT posts.post_id, people.username,
	title, time_creation FROM posts INNER JOIN people ON posts.user_id = people.user_id;`)
	CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		p := Post{}
		var t time.Time
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.TimeCreation)
		CheckErr(err)
		p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		allPosts = append(allPosts, p)
	}
	err = tmpl.ExecuteTemplate(w, "index.html", allPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AppContext) GetUser(id int) string {
	rows, err := c.db.Query(`SELECT username FROM people WHERE id = ?`, id)
	CheckErr(err)
	defer rows.Close()
	var username string
	for rows.Next() {
		err := rows.Scan(&username)
		CheckErr(err)
	}
	return username
}

func (c *AppContext) newCategory(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		title := r.FormValue("category")

		smt, err := c.db.Prepare(`INSERT INTO categories (title) VALUES(?)`)
		_, err = smt.Exec(title)
		CheckErr(err)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := tmpl.ExecuteTemplate(w, "newcategory.html", nil)
	CheckErr(err)
}

func (c *AppContext) comment(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("session")
		CheckErr(err)
		userID := c.GetSessions()[cookie.Value]
		postID := r.FormValue("postID")
		content := r.FormValue("content")

		stmt, err := c.db.Prepare(`INSERT INTO comments(user_id, 
			post_id, content, time_creation) VALUES (?, ?, ?, ?)`)
		CheckErr(err)
		_, err = stmt.Exec(userID, postID, content, time.Now())
		CheckErr(err)

		url := "/post/" + postID
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

}

func (c *AppContext) alreadyLogIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := c.GetSessions()[cookie.Value]

	return ok
}

// getSession gets session from database
func (c *AppContext) GetSessions() map[string]int {
	mapSessions := map[string]int{}

	rows, err := c.db.Query(`SELECT user_id, session_id FROM sessions`)
	CheckErr(err)

	for rows.Next() {
		var (
			sessionID string
			userID    int
		)
		err := rows.Scan(&userID, &sessionID)
		CheckErr(err)
		mapSessions[sessionID] = userID
	}
	return mapSessions
}
