package server

import (
	"forum/internal"
	"forum/pkg/models"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		s.badReq(w)
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
		s.methodNotAllowed(w)
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
		err = s.Template.ExecuteTemplate(w, "signin.html", models.Err{ErrCode: 401, ErrMsg: "invalid login"})
		if err != nil {
			log.Println(err)
			s.serverErr(w)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(clientPass))
	if err != nil {
		w.WriteHeader(403)
		err := s.Template.ExecuteTemplate(w, "signin.html", models.Err{ErrCode: 403, ErrMsg: "invalid password"})
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

	s.InfoLog.Println(u.Email, "signed-in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
