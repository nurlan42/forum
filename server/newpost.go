package server

import (
	"net/http"
	"strconv"
	"time"
)

func (c *AppContext) newPost(w http.ResponseWriter, r *http.Request) {

	if !c.alreadyLogIn(r) {
		ErrorPage(w, http.StatusForbidden, "please log-in first")
		return
	}

	// retrieving categoris to display user
	rows, err := c.db.Query(`SELECT * FROM categories`)
	CheckErr(err)
	defer rows.Close()

	categories := map[string]int{}
	for rows.Next() {
		var (
			id    int
			title string
		)

		err := rows.Scan(&id, &title)
		CheckErr(err)

		categories[title] = id

	}

	// recieving data from user
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("session")
		CheckErr(err)
		userID := c.GetSessions()[cookie.Value]
		title := r.FormValue("title")
		content := r.FormValue("post")
		// putting recieved data into database
		res, err := c.db.Exec(`INSERT INTO posts (user_id, title, content,
			 time_creation) VALUES(?, ?, ?, ?)`, userID, title, content, time.Now())
		CheckErr(err)
		postID, err := res.LastInsertId()
		CheckErr(err)

		stmt, err := c.db.Prepare(`INSERT INTO post_category(post_id, category_id) VALUES(?, ?);`)
		CheckErr(err)
		categories := r.Form["category"]

		for _, el := range categories {
			categoryID, err := strconv.Atoi(string(el))
			CheckErr(err)
			stmt.Exec(postID, categoryID)
		}
		stmt.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tmpl.ExecuteTemplate(w, "newpost.html", categories)
	CheckErr(err)
}
