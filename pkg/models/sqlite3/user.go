package sqlite3

import (
	"database/sql"
	"time"

	"forum/pkg/models"
)

func (c *Database) GetUser(uEmail string) (*models.User, error) {
	var u models.User
	row := c.SQLDb.QueryRow("SELECT user_id, email, password FROM people WHERE email = ?;", uEmail)
	err := row.Scan(&u.UserID, &u.Email, &u.Password)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	return &u, nil
}

// checking email for uniqness
func (c *Database) HasEmail(email string) bool {
	row := c.SQLDb.QueryRow(`SELECT email FROM people WHERE email = ?;`, email)
	err := row.Scan()
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (c *Database) InsertUser(u *models.User) (int64, error) {
	stmt, err := c.SQLDb.Prepare("INSERT INTO people (email, username, password, time_creation) VALUES(?, ?, ?, ?)")
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
