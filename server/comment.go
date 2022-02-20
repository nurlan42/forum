package server

import (
	"net/http"
	"strconv"
)

func (s *AppContext) commentNew(w http.ResponseWriter, r *http.Request) {
	if !s.alreadyLogIn(r) {
		s.ErrorHandler(w, http.StatusForbidden, "please, log-in first")
		return
	}
	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	userID, _ := s.Sqlite3.GetUserID(cookie.Value)
	s.Sqlite3.UpdateSession(userID)

	if r.Method != http.MethodPost {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Not Allowed")
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
