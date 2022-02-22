package server

import (
	"forum/internal"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	ok := s.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		err := s.Template.ExecuteTemplate(w, "signin.html", nil)
		//if template does not exist
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		s.signinPost(w, r)
	default:
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}

func (s *AppContext) signinPost(w http.ResponseWriter, r *http.Request) {
	var clientEmail, clientPass string

	clientEmail = r.FormValue("uemail")
	clientPass = r.FormValue("upass")

	// getting data from database, and saving into the var
	u, err := s.Sqlite3.GetUser(clientEmail)
	// err for incorrect login
	if err != nil {
		errorMsg := struct {
			Msg string
		}{
			"incorrect login",
		}
		// 401 unauthorised
		w.WriteHeader(401)
		err = s.Template.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(clientPass))
	if err != nil {
		errorMsg := struct {
			Msg string
		}{
			"incorrect password",
		}
		w.WriteHeader(403)
		err := s.Template.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	if s.Sqlite3.HasSession(u.UserID) {
		s.Sqlite3.DeleteSession(u.UserID)
	}
	//create new function
	sID := internal.SetCookie(w)

	s.Sqlite3.InsertSession(u.UserID, sID.String())

	http.Redirect(w, r, "/", http.StatusSeeOther)
	s.InfoLog.Println(u.Email, "signed-up successfully")

}
