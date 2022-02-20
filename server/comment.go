package server

import (
	"net/http"
	"strconv"
)

func (s *AppContext) comment(w http.ResponseWriter, r *http.Request) {
	if !s.alreadyLogIn(r) {
		s.ErrorHandler(w, http.StatusForbidden, "please, log-in first")
		return
	}
	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	userID, err := s.Sqlite3.GetUserID(cookie.Value)
	if err != nil {
		s.ErrorHandler(w, 500, "Internal Server Error")
		return
	}
	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := s.Sqlite3.GetUserID(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
			return
		}
		content := r.FormValue("content")
		err = s.Sqlite3.AddComment(userID, postID, content)
		if err != nil {
			s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		url := "/post/" + strconv.Itoa(postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

}
