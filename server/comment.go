package server

import (
	"net/http"
	"time"
)

func (c *AppContext) comment(w http.ResponseWriter, r *http.Request) {
	if !c.alreadyLogIn(r) {
		ErrorHandler(w, http.StatusForbidden, "please, log-in first")
		return
	}
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
