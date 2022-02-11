package sqlite3

import (
	"forum/pkg/models"
	"net/http"
	"strconv"
	"time"
)

func (c *Database) CreatePost(p *models.Post) (int64, error) {
	res, err := c.SqlDb.Exec(`INSERT INTO posts (user_id, title, content,
		time_creation) VALUES(?, ?, ?, ?)`, p.UserID, p.Title, p.Content, time.Now())
	if err != nil {
		return 0, err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postID, nil
}

// CreatePostCategory
func (c *Database) AddPostCategory(r *http.Request, postID int64) error {
	stmt, err := c.SqlDb.Prepare(`INSERT INTO post_category(post_id, category_id) VALUES(?, ?);`)
	if err != nil {
		return err
	}
	categories := r.Form["category"]

	for _, el := range categories {
		categoryID, err := strconv.Atoi(string(el))
		if err != nil {
			return err
		}
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			return err
		}
	}
	stmt.Close()
	return nil
}

// GetAllPosts gets all the posts from db
func (c *Database) GetAllPosts() (*[]models.Post, error) {
	var allPosts []models.Post
	rows, err := c.SqlDb.Query(`SELECT posts.post_id, people.username,
	title, posts.time_creation FROM posts INNER JOIN people ON posts.user_id = people.user_id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := models.Post{}
		var t time.Time
		err := rows.Scan(&p.PostID, &p.Author, &p.Title, &t)
		if err != nil {
			return nil, err
		}
		p.TimeCreation = t.Format("01-02-2006 15:04:05")
		p.CommentNbr, err = c.GetCommentsNbr(p.PostID)
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, p)
	}
	return &allPosts, nil
}

func (c *Database) GetPostByPostID(postID int) (*models.Post, error) {
	var p models.Post
	var t time.Time

	// putting data into struct from db
	row := c.SqlDb.QueryRow(`SELECT posts.post_id, people.username, title, content, 
	posts.time_creation FROM posts INNER JOIN people on posts.user_id = people.user_id WHERE post_id = ?`, postID)
	err := row.Scan(&p.PostID, &p.Author, &p.Title, &p.Content, &t)
	if err != nil {
		return nil, err
	}

	p.TimeCreation = t.Format("01-02-2006 15:04:05")

	return &p, nil
}
