package server

import "net/http"

func (s *AppContext) categoryNew(w http.ResponseWriter, r *http.Request) {

	if ok := s.alreadyLogIn(r); !ok {
		s.ErrorHandler(w, http.StatusForbidden, "please, log-in first")
		return
	}

	if r.URL.Path != "/category/new" {
		s.ErrorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	userID, err := s.Sqlite3.GetUserID(cookie.Value)
	CheckErr(err)

	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	switch r.Method {
	case http.MethodGet:
		err = s.Template.ExecuteTemplate(w, "newcategory.html", nil)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
	case http.MethodPost:
		title := r.FormValue("category")

		err = s.Sqlite3.InsertCategory(title)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}
