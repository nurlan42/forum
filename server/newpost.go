package server

import (
	"forum/pkg/models"
	"log"
	"net/http"
)

func (s *AppContext) showNewPost(w http.ResponseWriter, r *http.Request) {

	if !s.alreadyLogIn(r) {
		s.ErrorHandler(w, http.StatusForbidden, "please log-in first")
		return
	}

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300 // 300 is session length

	// update session table last activity
	mapSessID, err := s.Sqlite3.GetSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := mapSessID[cookie.Value]
	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	if r.URL.Path != "/newpost" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	switch r.Method {
	case http.MethodGet:
		categories, err := s.Sqlite3.GetAllCategories()
		if err != nil {
			s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = s.Template.ExecuteTemplate(w, "newpost.html", categories)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		s.newPost(w, r)
	default:
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")

	}

}

func (s *AppContext) newPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	CheckErr(err)

	session, err := s.Sqlite3.GetSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var p models.Post
	p.UserID = session[cookie.Value]
	p.Title = r.FormValue("title")
	p.Content = r.FormValue("post")

	// putting recieved data into database
	postID, err := s.Sqlite3.CreatePost(&p)
	if err != nil {
		if err != nil {
			s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	err = s.Sqlite3.AddPostCategory(r, postID)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
