package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.01.alem.school/Nurlan/forum.git/server/internal"
	"golang.org/x/crypto/bcrypt"
)

func (c *AppContext) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	ok := c.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		err := tmpl.ExecuteTemplate(w, "login.html", nil)
		//if template does not exist
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		c.loginPost(w, r)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}

// GetUser gets data from DB people
func (c *AppContext) GetUser(uEmail string) (*User, error) {
	var u User
	row := c.db.QueryRow("SELECT user_id, email, password FROM people WHERE email = ?;", uEmail)
	err := row.Scan(&u.ID, &u.Email, &u.Pass)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	return &u, nil

}

func (c *AppContext) loginPost(w http.ResponseWriter, r *http.Request) {
	var uEmail, uPass string

	uEmail = r.FormValue("uemail")
	uPass = r.FormValue("upass")

	// getting data from database, and saving into the var
	u, err := c.GetUser(uEmail)
	if err != nil {
		errorMsg := struct {
			Msg   string
			Email string
		}{
			"incorrect login",
			"",
		}
		// 401 unauthorised
		w.WriteHeader(401)
		err = tmpl.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Pass, []byte(uPass))
	if err != nil {
		errorMsg := struct {
			Msg   string
			Email string
		}{
			"incorrect password",
			uEmail,
		}
		w.WriteHeader(403)
		err := tmpl.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	if c.hasSession(u.ID) {
		c.DeleteSession(u.ID)
	}
	//create new function
	sID := internal.SetCookie(w)

	c.writeSession(u.ID, sID.String())

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("========== Logged-in successfully ==========")

}

func (c *AppContext) writeSession(userID int, sID string) {
	stmt, err := c.db.Prepare(`INSERT INTO sessions(user_id, session_id, start_date, expire_date) VALUES(?, ?, ?, ?);`)
	CheckErr(err)
	t := time.Now()

	stmt.Exec(userID, sID, t, t.Add(time.Minute*10))

}

func (c *AppContext) hasSession(userID int) bool {
	row := c.db.QueryRow(`SELECT session_id FROM sessions WHERE user_id = ?;`, userID)
	err := row.Scan()

	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}
