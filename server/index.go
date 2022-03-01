package server

import (
	"forum/pkg/models"
	"net/http"
)

func (s *AppContext) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.ErrorLog.Println(r.URL.Path)
		s.badReq(w)
		return
	}

	if r.Method != http.MethodGet {
		s.methodNotAllowed(w)
		return
	}

	// new function
	allPosts, err := s.Sqlite3.GetAllPosts()
	if err != nil {
		s.serverErr(w)
		return
	}

	// categories
	categories, err := s.Sqlite3.GetAllCategories()
	if err != nil {
		s.serverErr(w)
		return
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
