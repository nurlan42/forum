package sqlite3

import "database/sql"

// writePostLike adds 1 to post_reaction
func (c *Database) AddPostReaction(userID, postID, reaction int) error {
	_, err := c.SqlDb.Exec(`INSERT INTO post_reaction(user_id, post_id, reaction) VALUES(?, ?, ?)`, userID, postID, reaction)
	if err != nil {
		return err
	}
	return nil
}
func (c *Database) UpdatePostReaction(userID, postID, reaction int) error {
	_, err := c.SqlDb.Exec(`UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?;`, reaction, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) DeletePostReaction(userID, postID int) error {
	_, err := c.SqlDb.Exec(`DELETE FROM post_reaction WHERE user_id = ? AND post_id = ?;`, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) AddCommReaction(userID, commID, reaction int) error {
	_, err := c.SqlDb.Exec(`INSERT INTO comment_reaction(user_id, comment_id, reaction) VALUES(?, ?, ?)`, userID, commID, reaction)
	if err != nil {
		return err
	}
	return nil
}
func (c *Database) UpdateCommReaction(userID, commID, reaction int) error {
	_, err := c.SqlDb.Exec(`UPDATE comment_reaction SET reaction = ? WHERE user_id = ? AND comment_id = ?;`, reaction, userID, commID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) DeleteCommReaction(userID, commID int) error {
	_, err := c.SqlDb.Exec(`DELETE FROM comment_reaction WHERE user_id = ? AND comment_id = ?;`, userID, commID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) HasReactionPost(userID, postID int) (bool, int) {
	var r int
	row := c.SqlDb.QueryRow(`SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?;`, userID, postID)
	err := row.Scan(&r)
	if err == sql.ErrNoRows {
		return false, -1
	}
	return true, r
}

func (c *Database) HasReactionComm(userID, commID int) (bool, int) {
	var r int
	row := c.SqlDb.QueryRow(`SELECT reaction FROM comment_reaction WHERE user_id = ? AND comment_id = ?;`, userID, commID)
	err := row.Scan(&r)
	if err == sql.ErrNoRows {
		return false, -1
	}
	return true, r
}

func (c *Database) ReadPostID(commID int) (int, error) {
	var postID int
	row := c.SqlDb.QueryRow(`SELECT post_id FROM comments WHERE comment_id = ?;`, commID)
	err := row.Scan(&postID)
	if err == sql.ErrNoRows {
		return -1, err
	}
	return postID, nil

}

func (c *Database) GetPostReaction(postID int) (int, int, error) {
	var like, dislike int
	row := c.SqlDb.QueryRow(`SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = ?;`, postID, 1)
	err := row.Scan(&like)
	if err != nil {
		return 0, 0, err
	}

	row = c.SqlDb.QueryRow(`SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = ?;`, postID, 0)
	err = row.Scan(&dislike)
	if err != nil {
		return 0, 0, err
	}
	return like, dislike, nil
}
