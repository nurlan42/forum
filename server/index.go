package server

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("ui/html/*.html"))

}

func (c *AppContext) index(w http.ResponseWriter, r *http.Request) {

	if !c.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// new function
	allPosts, err := c.ReadPosts()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", allPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AppContext) ReadPosts() (*[]Post, error) {
	var allPosts []Post
	rows, err := c.db.Query(`SELECT posts.post_id, people.username,
	title, time_creation FROM posts INNER JOIN people ON posts.user_id = people.user_id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Post{}
		var t time.Time
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.TimeCreation)
		CheckErr(err)
		p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		allPosts = append(allPosts, p)
	}
	return &allPosts, nil
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

		session, err := c.getSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session[cookie.Value]
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
	// _, ok := c.GetSessions()[cookie.Value]
	_, err = c.getSession(cookie.Value)
	if err != nil {
		return false
	}

	return true
}

func (c *AppContext) getSession(s string) (map[string]int, error) {
	mapSession := map[string]int{}

	var (
		sessionID string
		userID    int
	)

	row := c.db.QueryRow(`SELECT user_id, session_id FROM sessions WHERE session_id = ?;`, s)
	err := row.Scan(&userID, &sessionID)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	mapSession[sessionID] = userID
	return mapSession, nil
}
