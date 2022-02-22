package sqlite3

func (c *Database) InsertCategory(title string) error {
	smt, err := c.SQLDb.Prepare(`INSERT INTO categories (title) VALUES(?)`)
	_, err = smt.Exec(title)
	if err != nil {
		return err
	}
	return nil
}

func (c *Database) GetCategoriesByPostID(postID int) ([]string, error) {
	var categories []string
	rows, err := c.SQLDb.Query(`SELECT categories.title FROM post_category INNER JOIN categories 
	on post_category.category_id = categories.category_id WHERE post_category.post_id = ?;`, postID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			return []string{}, err
		}
		categories = append(categories, title)
	}
	return categories, nil
}

func (c *Database) GetAllCategories() (map[string]int, error) {
	rows, err := c.SQLDb.Query(`SELECT * FROM categories`)
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
