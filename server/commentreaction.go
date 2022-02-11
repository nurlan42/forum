package server

import (
	"net/http"
	"strconv"
)

func (s *AppContext) commentReaction(w http.ResponseWriter, r *http.Request) {
	if !s.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	cookie, _ := r.Cookie("session")
	commID, err := strconv.Atoi(r.FormValue("commID"))
	CheckErr(err)

	mapSessID, _ := s.Sqlite3.GetSession(cookie.Value)
	userID := mapSessID[cookie.Value]

	// 0 is dislike, 1 is like
	reaction, err := strconv.Atoi(r.FormValue("reaction"))
	CheckErr(err)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}

	b, e := s.Sqlite3.HasReactionComm(userID, commID)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if b {
		if e == reaction {
			s.Sqlite3.DeleteCommReaction(userID, commID)
		} else {
			s.Sqlite3.UpdateCommReaction(userID, commID, reaction)
		}
	} else {
		s.Sqlite3.AddCommReaction(userID, commID, reaction)
	}
	postID, err := s.Sqlite3.ReadPostID(commID)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	url := "/post/" + strconv.Itoa(postID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
