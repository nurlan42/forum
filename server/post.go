package server

import (
	"fmt"
	"forum/internal"
	"log"
	"net/http"
)

// showPost handler
func (s *AppContext) post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	//get id of post new func
	postID, err := internal.GetPostID(r)
	if err != nil {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	// query to get the post from db
	p, err := s.Sqlite3.GetPostByPostID(postID)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Not Found")
		return
	}

	// retrieve categories from db
	p.Categories, err = s.Sqlite3.GetCategoriesByPostID(postID)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Not Found")
		return
	}

	// get reaction nbr for a post
	p.Like, p.Dislike, err = s.Sqlite3.GetPostReaction(postID)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// retrive comments from db
	p.Comments, err = s.Sqlite3.GetCommentsByPostID(postID)
	if err != nil {
		fmt.Println(err)
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = s.Template.ExecuteTemplate(w, "post.html", p)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
