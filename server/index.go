package server

import (
	"forum/pkg/models"
	"net/http"
)

func (s *AppContext) index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if r.Method != http.MethodGet {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// new function
	allPosts, err := s.Sqlite3.GetAllPosts()
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// categories
	categories, err := s.Sqlite3.GetAllCategories()
	CheckErr(err)

	data := struct {
		AllPosts   *[]models.Post
		Categories map[string]int
	}{allPosts, categories}

	err = s.Template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
