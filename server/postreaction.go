package server

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *AppContext) postReaction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, _ := r.Cookie("session")
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/reaction/"))

	if err != nil {
		s.badReq(w)
		return
	}

	userID, _ := s.Sqlite3.GetUserID(cookie.Value)

	// 0 is dislike, 1 is like
	reaction, err := strconv.Atoi(r.FormValue("reaction"))
	if err != nil {
		s.badReq(w)
		return
	}

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}

	ok, emotion := s.Sqlite3.HasReactionPost(userID, postID)

	if ok {
		if emotion == reaction {
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
