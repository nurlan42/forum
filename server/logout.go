package server

import (
	"fmt"
	"net/http"
)

func (c *AppContext) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	ok := c.alreadyLogIn(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("cookie unexist")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// mapSession stores [session] = userID
	mapSession, err := c.getSession(cookie.Value)
	CheckErr(err)

	userID := mapSession[cookie.Value]
	c.DeleteSession(userID)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AppContext) DeleteSession(userID int) {
	stmt, err := c.db.Prepare(`DELETE FROM sessions 
		WHERE user_id = ?;`)
	CheckErr(err)
	stmt.Exec(userID)
	stmt.Close()
}
