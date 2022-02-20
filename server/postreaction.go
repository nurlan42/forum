package server

import (
	"net/http"
	"strconv"
)

func (s *AppContext) postReaction(w http.ResponseWriter, r *http.Request) {

	if !s.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	cookie, _ := r.Cookie("session")
	postID, err := strconv.Atoi(r.FormValue("postID"))
	CheckErr(err)

	userID, _ := s.Sqlite3.GetUserID(cookie.Value)

	// 0 is dislike, 1 is like
	reaction, err := strconv.Atoi(r.FormValue("reaction"))
	CheckErr(err)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}

	b, e := s.Sqlite3.HasReactionPost(userID, postID)

	if b {
		if e == reaction {
			s.Sqlite3.DeletePostReaction(userID, postID)
		} else {
			s.Sqlite3.UpdatePostReaction(userID, postID, reaction)
		}
	} else {
		s.Sqlite3.AddPostReaction(userID, postID, reaction)
	}

	url := "/post/" + strconv.Itoa(postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
