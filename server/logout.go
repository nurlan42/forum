package server

import (
	"fmt"
	"net/http"
)

func (s *AppContext) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	ok := s.alreadyLogIn(r)
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
	mapSession, err := s.Sqlite3.GetSession(cookie.Value)
	CheckErr(err)

	userID := mapSession[cookie.Value]
	s.Sqlite3.DeleteSession(userID)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
