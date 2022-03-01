package server

import (
	"forum/internal"
	"net/http"
)

// showPost handler
func (s *AppContext) post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.methodNotAllowed(w)
		return
	}

	//get id of post new func
	postID, err := internal.GetPostID(r)
	if err != nil {
		s.badReq(w)
		return
	}
	// query to get the post from db
	p, err := s.Sqlite3.GetPostByPostID(postID)
	if err != nil {
		s.notFound(w)
		return
	}

	// retrieve categories from db
	p.Categories, err = s.Sqlite3.GetCategoriesByPostID(postID)
	if err != nil {
		s.notFound(w)
		return
	}

	// get reaction nbr for a post
	p.Like, p.Dislike, err = s.Sqlite3.GetPostReaction(postID)
	if err != nil {
		s.notFound(w)
		return
	}

	// retrive comments from db
	p.Comments, err = s.Sqlite3.GetCommentsByPostID(postID)
	if err != nil {
		s.ErrorLog.Println(err)
		s.notFound(w)
		return
	}

	err = s.Template.ExecuteTemplate(w, "post.html", p)
	if err != nil {
		s.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
