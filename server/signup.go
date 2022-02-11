package server

import (
	"forum/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	ok := s.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		var u models.User
		u.Email = r.FormValue("uemail")

		if s.Sqlite3.HasEmail(u.Email) {
			s.ErrorHandler(w, http.StatusNotAcceptable, "That email already occupied, Try another.")
			return
		}

		u.UserName = r.FormValue("uname")
		pass := r.FormValue("upass")

		var err error
		u.Password, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		CheckErr(err)

		_, err = s.Sqlite3.InsertUser(&u)
		CheckErr(err)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")

}

func (s *AppContext) alreadyLogIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// getSession from db
	_, err = s.Sqlite3.GetSession(cookie.Value)

	return err == nil
}
