package sqlite3

import (
	"forum/pkg/models"
	"time"
)

func (c *Database) AddComment(userID, postID int, content string) error {
	stmt, err := c.SqlDb.Prepare(`INSERT INTO comments(user_id, 
		post_id, content, time_creation) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, postID, content, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	rows, err := c.SqlDb.Query(`SELECT comments.comment_id, people.username, content, 
	comments.time_creation FROM comments INNER JOIN PEOPLE on comments.user_id = people.user_id WHERE post_id = ?`, postID)
	if err != nil {
		return nil, err
	}

	var comments []models.Comment
	for rows.Next() {
		var t time.Time
		var comment models.Comment
		err := rows.Scan(&comment.CommID, &comment.Author, &comment.Content, &t)
		if err != nil {
			return nil, err
		}
		comment.Like, comment.Dislike, err = c.GetCommentReaction(comment.CommID)
		if err != nil {
			return nil, err
		}
		comment.TimeCreation = t.Format("01-02-2006 15:04:05 Monday")
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *Database) GetCommentReaction(commID int) (int, int, error) {
	var like, dislike int
	row := c.SqlDb.QueryRow(`SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = ?;`, commID, 1)
	err := row.Scan(&like)
	if err != nil {
		return 0, 0, err
	}

	row = c.SqlDb.QueryRow(`SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = ?;`, commID, 0)
	err = row.Scan(&dislike)
	if err != nil {
		return 0, 0, err
	}
	return like, dislike, nil
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
