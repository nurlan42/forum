package server

import (
	"forum/pkg/models"
	"net/http"
	"strconv"
)

func (s *AppContext) filter(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/filter" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if r.Method != http.MethodPost {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	allPosts, err := s.Sqlite3.GetAllPosts()
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	categories, err := s.Sqlite3.GetAllCategories()
	CheckErr(err)

	if r.FormValue("owner") == "yes" {
		allPosts, err = s.filterByOwner(r)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
	} else if r.FormValue("category") != "" {
		categoryID, err := strconv.Atoi(r.FormValue("category"))
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
		allPosts, err = s.Sqlite3.GetPostsByCategory(categoryID)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
	} else if r.FormValue("reaction") != "" {
		reaction, err := strconv.Atoi(r.FormValue("reaction"))
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
		allPosts, err = s.Sqlite3.GetPostsByReaction(reaction)
		if err != nil {
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
	}

	data := struct {
		AllPosts   *[]models.Post
		Categories map[string]int
	}{allPosts, categories}

	err = s.Template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *AppContext) filterByOwner(r *http.Request) (*[]models.Post, error) {
	// created by you
	cookie, err := r.Cookie("session")
	CheckErr(err)

	sID, err := s.Sqlite3.GetSession(cookie.Value)
	CheckErr(err)

	userID := sID[cookie.Value]

	ps, err := s.Sqlite3.GetPostsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
