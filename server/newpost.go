package server

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func (c *AppContext) showNewPost(w http.ResponseWriter, r *http.Request) {

	if !c.alreadyLogIn(r) {
		ErrorHandler(w, http.StatusForbidden, "please log-in first")
		return
	}

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300 // 300 is session length

	// update session table last activity
	mapSessID, err := c.getSession(cookie.Value)
	CheckErr(err)
	userID := mapSessID[cookie.Value]
	if c.hasSession(userID) {
		c.updateSession(userID)
	}

	if r.URL.Path != "/newpost" {
		ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	switch r.Method {
	case http.MethodGet:
		categories, err := c.readCategories()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = tmpl.ExecuteTemplate(w, "newpost.html", categories)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		c.newPost(w, r)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")

	}

}

func (c *AppContext) readCategories() (map[string]int, error) {
	rows, err := c.db.Query(`SELECT * FROM categories`)
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

func (c *AppContext) newPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	CheckErr(err)

	session, err := c.getSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var p Post
	p.UserID = session[cookie.Value]
	p.Title = r.FormValue("title")
	p.Content = r.FormValue("post")

	// putting recieved data into database
	postID, err := c.writePost(&p)
	if err != nil {
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	err = c.writePostCategory(r, postID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (c *AppContext) writePost(p *Post) (int64, error) {
	res, err := c.db.Exec(`INSERT INTO posts (user_id, title, content,
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

func (c *AppContext) writePostCategory(r *http.Request, postID int64) error {
	stmt, err := c.db.Prepare(`INSERT INTO post_category(post_id, category_id) VALUES(?, ?);`)
	CheckErr(err)
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
