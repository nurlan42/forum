package server

import "net/http"

func (s *AppContext) newCategory(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	mapSessID, err := s.Sqlite3.GetSession(cookie.Value)
	CheckErr(err)
	userID := mapSessID[cookie.Value]
	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("category")

		err = s.Sqlite3.AddCategory(title)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = s.Template.ExecuteTemplate(w, "newcategory.html", nil)
	if err != nil {
		s.ErrorHandler(w, 500, "Internal Server Error")
		return
	}
}
