package server

import (
	"net/http"
)

func (s *AppContext) auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		_, err = s.Sqlite3.GetUserID(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		HandlerFunc.ServeHTTP(w, r)
	}

}
