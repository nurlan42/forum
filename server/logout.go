package server

import (
	"fmt"
	"net/http"
)

func (c *AppContext) logout(w http.ResponseWriter, r *http.Request) {
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

	userID := c.GetSessions()[cookie.Value]

	smt, _ := c.db.Prepare(`DELETE FROM sessions 
		WHERE user_id = ?`)
	smt.Exec(userID)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
