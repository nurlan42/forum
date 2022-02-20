package sqlite3

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Database) InsertSession(userID int, sID string) error {
	stmt, err := c.SqlDb.Prepare(`INSERT INTO sessions(user_id, session_id, start_date, expire_date) VALUES(?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	t := time.Now()
	stmt.Exec(userID, sID, t, t.Add(time.Minute*10))
	return nil
}
func (c *Database) UpdateSession(userID int) error {
	_, err := c.SqlDb.Exec(`UPDATE sessions SET start_date = ?, expire_date = ? WHERE user_id = ?;`, time.Now(), time.Now().Add(time.Minute*5), userID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) DeleteInactiveSession() error {
	fmt.Println("DeleteInactiveSession()")
	_, err := c.SqlDb.Exec(`DELETE from sessions WHERE expire_date <= ?;`, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) DeleteSession(userID int) error {
	stmt, err := c.SqlDb.Prepare(`DELETE FROM sessions 
		WHERE user_id = ?;`)
	if err != nil {
		return err
	}
	stmt.Exec(userID)
	stmt.Close()
	return nil
}

func (c *Database) HasSession(userID int) bool {
	row := c.SqlDb.QueryRow(`SELECT session_id FROM sessions WHERE user_id = ?;`, userID)
	err := row.Scan()

	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

// getSession gets session from db
func (c *Database) GetUserID(sID string) (int, error) {
	var userID int

	row := c.SqlDb.QueryRow(`SELECT user_id FROM sessions WHERE session_id = ?;`, sID)
	err := row.Scan(&userID)
	if err != nil && err == sql.ErrNoRows {
		return 0, err
	}

	return userID, nil
}
