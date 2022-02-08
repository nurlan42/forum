package server

import (
	"log"
	"net/http"
	"time"

	"git.01.alem.school/Nurlan/forum.git/server/internal"
)

func (c *AppContext) readPosts(postID int) (*Post, error) {
	var p Post
	var t time.Time

	// putting data into struct from db
	row := c.db.QueryRow(`SELECT posts.post_id, people.username, title, content, 
	time_creation FROM posts INNER JOIN people on posts.user_id = people.user_id WHERE post_id = ?`, postID)
	err := row.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &t)
	if err != nil {
		return nil, err
	}

	p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")

	return &p, nil
}

// readPostReaction gets reaction nbr for a post
func (c *AppContext) readPostReaction(postID int) (like, dislike int) {
	// var like, dislike int
	row := c.db.QueryRow(`SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = ?;`, postID, 1)
	err := row.Scan(&like)
	CheckErr(err)

	row = c.db.QueryRow(`SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = ?;`, postID, 0)
	err = row.Scan(&dislike)
	CheckErr(err)
	return like, dislike
}

func (c *AppContext) ReadCategories(postID int) ([]string, error) {
	var categories []string
	rows, err := c.db.Query(`SELECT categories.title FROM post_category INNER JOIN categories 
	on post_category.category_id = categories.category_id WHERE post_category.post_id = ?;`, postID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		CheckErr(err)
		categories = append(categories, title)
	}
	return categories, nil
}

func (c *AppContext) ReadComments(postID int) ([]Comment, error) {
	rows, err := c.db.Query(`SELECT comments.comment_id, people.username, content, 
	time_creation FROM comments INNER JOIN PEOPLE on comments.user_id = people.user_id WHERE post_id = ?`, postID)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	for rows.Next() {
		var t time.Time
		var comment Comment
		err := rows.Scan(&comment.CommID, &comment.Author, &comment.Content, &t)
		if err != nil {
			return nil, err
		}
		comment.Like, comment.Dislike = c.readCommReaction(comment.CommID)
		comment.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *AppContext) readCommReaction(commID int) (like, dislike int) {
	row := c.db.QueryRow(`SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = ?;`, commID, 1)
	err := row.Scan(&like)
	CheckErr(err)

	row = c.db.QueryRow(`SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = ?;`, commID, 0)
	err = row.Scan(&dislike)
	CheckErr(err)
	return
}

// showPost handler
func (c *AppContext) showPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	//get id of post new func
	postID, err := internal.GetPostID(r)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	// query to get the post from db
	p, err := c.readPosts(postID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Not Found")
		return
	}

	// retrieve categories from db
	p.Categories, err = c.ReadCategories(postID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Not Found")
		return
	}

	// get reaction nbr for a post
	p.Like, p.Dislike = c.readPostReaction(postID)

	// retrive comments from db
	p.Comments, err = c.ReadComments(postID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = tmpl.ExecuteTemplate(w, "post.html", p)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
