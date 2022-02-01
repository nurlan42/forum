package server

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (c *AppContext) signup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	ok := c.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("uemail")

		if c.hasEmail(email) {
			ErrorHandler(w, http.StatusNotAcceptable, "That email already occupied, Try another.")
			return
		}

		userName := r.FormValue("uname")
		pass := r.FormValue("upass")

		passBs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		CheckErr(err)

		stmt, err :=
			c.db.Prepare("INSERT INTO people (email, username, password) VALUES(?, ?, ?)")
		_, err = stmt.Exec(email, userName, passBs)
		CheckErr(err)
		defer stmt.Close()

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if r.Method == http.MethodGet {
		err := tmpl.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
	}

}

// checking email for uniqness
func (c *AppContext) hasEmail(email string) bool {
	row := c.db.QueryRow(`SELECT email FROM people WHERE email = ?;`, email)
	err := row.Scan()
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true

}
