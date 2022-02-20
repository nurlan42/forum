package sqlite3

import (
	"forum/pkg/models"
	"net/http"
	"strconv"
	"time"
)

func (c *Database) InserPost(p *models.Post) (int64, error) {
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
func (c *Database) InsertPostCategory(r *http.Request, postID int64) error {
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

func (c *Database) GetPostsByUserID(userID int) (*[]models.Post, error) {

	rows, err := c.SqlDb.Query(`SELECT posts.post_id, people.username,title, content, 
	posts.time_creation FROM posts INNER JOIN people on posts.user_id = people.user_id WHERE posts.user_id = ?;`, userID)
	if err != nil {
		return nil, err
	}

	var ps []models.Post
	for rows.Next() {
		var p models.Post
		var t time.Time
		err := rows.Scan(&p.PostID, &p.Author, &p.Title, &p.Content, &t)
		if err != nil {
			return nil, err
		}
		p.TimeCreation = t.Format("01-02-2006 15:04:05")
		p.CommentNbr, err = c.GetCommentsNbr(p.PostID)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return &ps, nil
}

// GetPostsByCategory gets all the posts by category
func (c *Database) GetPostsByCategory(categoryID int) (*[]models.Post, error) {

	rows, err := c.SqlDb.Query(`SELECT posts.post_id, people.username, posts.title,
		posts.content, posts.time_creation FROM posts INNER JOIN people ON posts.user_id =
		people.user_id INNER JOIN post_category ON post_category.post_id = posts.post_id 
		WHERE post_category.category_id = ?`, categoryID)

	// rows, err := c.db.Query(`SELECT posts.post_id, posts.title, posts.content, posts.time_creation FROM
	// 	posts INNER JOIN post_category ON posts.post_id = post_category.post_id WHERE post_category.category_id = ?;`, categoryID)
	if err != nil {
		return nil, err
	}

	var ps []models.Post
	for rows.Next() {
		var p models.Post
		var t time.Time
		err := rows.Scan(&p.PostID, &p.Author, &p.Title, &p.Content, &t)
		if err != nil {
			return nil, err
		}
		p.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		p.CommentNbr, err = c.GetCommentsNbr(p.PostID)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return &ps, nil
}

func (c *Database) GetPostsByReaction(emotion, userID int) (*[]models.Post, error) {
	rows, err := c.SqlDb.Query(`SELECT posts.post_id, people.username, posts.title, posts.content, posts.time_creation FROM posts
		INNER JOIN people ON people.user_id = posts.user_id INNER JOIN post_reaction ON post_reaction.post_id = posts.post_id
		WHERE post_reaction.reaction = ? AND posts.user_id = ?;`, emotion, userID)
	if err != nil {
		return nil, err
	}

	var ps []models.Post
	for rows.Next() {
		var p models.Post
		var t time.Time
		err := rows.Scan(&p.PostID, &p.Author, &p.Title, &p.Content, &t)
		if err != nil {
			return nil, err
		}
		p.TimeCreation = t.Format("01-02-2006 15:04:05")
		p.CommentNbr, err = c.GetCommentsNbr(p.PostID)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return &ps, nil
}
