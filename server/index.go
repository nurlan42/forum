package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("ui/html/*.html"))

}

func (c *AppContext) index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// new function
	allPosts, err := c.ReadPosts()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// categories
	categories, err := c.readCategories()
	CheckErr(err)

	fmt.Println("commentNbr=", c.readCommentsNbr(1))

	data := struct {
		AllPosts   *[]Post
		Categories map[string]int
	}{allPosts, categories}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AppContext) readCommentsNbr(postID int) int {
	var i int
	row := c.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?;`, postID)
	err := row.Scan(&i)
	CheckErr(err)
	return i

}

// ReadPosts gets all the posts from db
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
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &t)
		CheckErr(err)
		p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		p.CommentNbr = c.readCommentsNbr(p.ID)
		allPosts = append(allPosts, p)
	}
	fmt.Println(allPosts)
	return &allPosts, nil
}

// newCategory is handles
func (c *AppContext) newCategory(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	mapSessID, err := c.getSession(cookie.Value)
	CheckErr(err)
	userID := mapSessID[cookie.Value]
	if c.hasSession(userID) {
		c.updateSession(userID)
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("category")

		smt, err := c.db.Prepare(`INSERT INTO categories (title) VALUES(?)`)
		_, err = smt.Exec(title)
		CheckErr(err)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tmpl.ExecuteTemplate(w, "newcategory.html", nil)
	CheckErr(err)
}

func (c *AppContext) alreadyLogIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// getSession from db
	_, err = c.getSession(cookie.Value)
	if err != nil {
		return false
	}

	return true
}

// getSession gets session from db
func (c *AppContext) getSession(sID string) (map[string]int, error) {
	mapSession := map[string]int{}

	var (
		sessionID string
		userID    int
	)

	row := c.db.QueryRow(`SELECT user_id, session_id FROM sessions WHERE session_id = ?;`, sID)
	err := row.Scan(&userID, &sessionID)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	mapSession[sessionID] = userID
	return mapSession, nil
}

func (c *AppContext) updateSession(userID int) {
	_, err := c.db.Exec(`UPDATE sessions SET start_date = ?, expire_date = ? WHERE user_id = ?;`, time.Now(), time.Now().Add(time.Minute*5), userID)
	CheckErr(err)
}

func (c *AppContext) DeleteInactiveSession() {
	fmt.Println("DeleteInactiveSession()")
	_, err := c.db.Exec(`DELETE from sessions WHERE expire_date <= ?;`, time.Now())
	CheckErr(err)
}
