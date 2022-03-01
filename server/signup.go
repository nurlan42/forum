package server

import (
	"forum/internal"
	"forum/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		s.badReq(w)
		return
	}

	if r.Method != http.MethodPost {
		s.methodNotAllowed(w)
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
		u.UserName = r.FormValue("uname")
		pass := r.FormValue("upass")
		passConfirm := r.FormValue("confirm-upass")

		if ok := internal.CheckName(u.UserName); !ok {
			s.clientErr(w, models.Err{ErrCode: http.StatusBadRequest, ErrMsg: "Please, be sure to write valid name"})
			return
		}
		if ok := internal.CheckEmail(u.Email); !ok {
			s.clientErr(w, models.Err{ErrCode: http.StatusBadRequest, ErrMsg: "Please, be sure to write valid email"})
			return
		}

		if pass != passConfirm {
			s.clientErr(w, models.Err{ErrCode: http.StatusBadRequest, ErrMsg: "password is empty / doesn't match"})
			return
		}

		if len(pass) < 5 {
			s.clientErr(w, models.Err{ErrCode: http.StatusBadRequest, ErrMsg: "Password should be more than 5 symbols"})
			return
		}

		var err error
		u.Password, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		if err != nil {
			s.ErrorLog.Println(err)
			s.serverErr(w)
			return
		}

		_, err = s.Sqlite3.InsertUser(&u)
		if err != nil {
			s.clientErr(w, models.Err{http.StatusNotAcceptable, "That email already occupied. Try another."})
			return
		}
		s.InfoLog.Println(u.Email, "signed-up successfully")
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

}

func (s *AppContext) alreadyLogIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// getSession from db
	_, err = s.Sqlite3.GetUserID(cookie.Value)
	if err != nil {
		return false
	}
	return err == nil
}
