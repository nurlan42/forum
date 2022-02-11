package sqlite3

func (c *Database) GetCategoriesByPostID(postID int) ([]string, error) {
	var categories []string
	rows, err := c.SqlDb.Query(`SELECT categories.title FROM post_category INNER JOIN categories 
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
