// Package pkg interacts with database directly
package sqlite3

import (
	"database/sql"
	"forum/pkg/models"
	"time"
)

type Database struct {
	SqlDb *sql.DB
}

func (c *Database) GetUser(uEmail string) (*models.User, error) {
	var u models.User
	row := c.SqlDb.QueryRow("SELECT user_id, email, password FROM people WHERE email = ?;", uEmail)
	err := row.Scan(&u.UserID, &u.Email, &u.Password)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	return &u, nil

}

// checking email for uniqness
func (c *Database) HasEmail(email string) bool {
	row := c.SqlDb.QueryRow(`SELECT email FROM people WHERE email = ?;`, email)
	err := row.Scan()
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true

}

func (c *Database) InsertUser(u *models.User) (int64, error) {
	stmt, err := c.SqlDb.Prepare("INSERT INTO people (email, username, password, time_creation) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	t := time.Now()
	res, err := stmt.Exec(u.Email, u.UserName, u.Password, t)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (c *Database) GetCommentsNbr(postID int) (int, error) {
	var i int
	row := c.SqlDb.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?;`, postID)
	err := row.Scan(&i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (c *Database) GetAllCategories() (map[string]int, error) {
	rows, err := c.SqlDb.Query(`SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := map[string]int{}
	for rows.Next() {
		var (
			id    int
			title string
		)

		err := rows.Scan(&id, &title)
		if err != nil {
			return nil, err
		}

		categories[title] = id
	}
	return categories, nil

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

func (c *Database) GetPostsByReaction(emotion int) (*[]models.Post, error) {
	rows, err := c.SqlDb.Query(`SELECT posts.post_id, people.username, posts.title, posts.content, posts.time_creation FROM posts
		INNER JOIN people ON people.user_id = posts.user_id INNER JOIN post_reaction ON post_reaction.post_id = posts.post_id
		WHERE post_reaction.reaction = ?;`, emotion)
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

func (c *Database) AddCategory(title string) error {
	smt, err := c.SqlDb.Prepare(`INSERT INTO categories (title) VALUES(?)`)
	_, err = smt.Exec(title)
	if err != nil {
		return err
	}
	return nil
}
