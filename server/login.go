package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (c *AppContext) login(w http.ResponseWriter, r *http.Request) {
	ok := c.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var uemail, upass string
	if r.Method == http.MethodPost {
		uemail = r.FormValue("uemail")
		upass = r.FormValue("upass")

		var u User // getting data from database, and saving into the var
		row := c.db.QueryRow("SELECT user_id, email, password FROM people WHERE email = ?;", uemail)
		err := row.Scan(&u.ID, &u.Email, &u.Pass)

		if err != nil && err == sql.ErrNoRows {
			ErrorPage(w, 403, "incorrect login")
			return
		}

		err = bcrypt.CompareHashAndPassword(u.Pass, []byte(upass))
		if err != nil {
			ErrorPage(w, 403, "incorrect password")
			return
		}

		sID := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, cookie)
	 
		// put data into database
		stmt, err := c.db.Prepare("INSERT INTO sessions(user_id, session_id, last_activity) VALUES(?, ?, ?)")
		CheckErr(err)
		t := time.Now()
		stmt.Exec(u.ID, sID, t)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("========== Logged-in successfully ==========")
		return

	}

	err := tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
