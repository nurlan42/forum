package server

import (
	"net/http"
)

func (s *AppContext) signout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signout" {
		s.badReq(w)
		return
	}

	ok := s.alreadyLogIn(r)
	if !ok {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	cookie, _ := r.Cookie("session")

	// mapSession stores [session] = userID
	userID, _ := s.Sqlite3.GetUserID(cookie.Value)

	s.Sqlite3.DeleteSession(userID)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
