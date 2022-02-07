package server

import (
	"database/sql"
	"net/http"
	"strconv"
)

func (c *AppContext) postReaction(w http.ResponseWriter, r *http.Request) {

	if !c.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	cookie, _ := r.Cookie("session")
	postID, err := strconv.Atoi(r.FormValue("postID"))
	CheckErr(err)

	mapSessID, _ := c.getSession(cookie.Value)
	userID := mapSessID[cookie.Value]

	// 0 is dislike, 1 is like
	reaction, err := strconv.Atoi(r.FormValue("reaction"))
	CheckErr(err)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}

	b, e := c.hasReactionPost(userID, postID)

	if b {
		if e == reaction {
			c.deletePostReaction(userID, postID)
		} else {
			c.updatePostReaction(userID, postID, reaction)
		}
	} else {
		c.writePostReaction(userID, postID, reaction)
	}

	url := "/post/" + strconv.Itoa(postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (c *AppContext) hasReactionPost(userID, postID int) (bool, int) {
	var r int
	row := c.db.QueryRow(`SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?;`, userID, postID)
	err := row.Scan(&r)
	if err == sql.ErrNoRows {
		return false, -1
	}
	return true, r
}

// writePostLike adds 1 to post_reaction
func (c *AppContext) writePostReaction(userID, postID, reaction int) {
	_, err := c.db.Exec(`INSERT INTO post_reaction(user_id, post_id, reaction) VALUES(?, ?, ?)`, userID, postID, reaction)
	CheckErr(err)
}
func (c *AppContext) updatePostReaction(userID, postID, reaction int) {
	_, err := c.db.Exec(`UPDATE post_reaction SET reaction = ? WHERE user_id = ? AND post_id = ?;`, reaction, userID, postID)
	CheckErr(err)
}

func (c *AppContext) deletePostReaction(userID, postID int) {
	_, err := c.db.Exec(`DELETE FROM post_reaction WHERE user_id = ? AND post_id = ?;`, userID, postID)
	CheckErr(err)
}

func (c *AppContext) commentReaction(w http.ResponseWriter, r *http.Request) {
	if !c.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	cookie, _ := r.Cookie("session")
	commID, err := strconv.Atoi(r.FormValue("commID"))
	CheckErr(err)

	mapSessID, _ := c.getSession(cookie.Value)
	userID := mapSessID[cookie.Value]

	// 0 is dislike, 1 is like
	reaction, err := strconv.Atoi(r.FormValue("reaction"))
	CheckErr(err)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}

	b, e := c.hasReactionComm(userID, commID)

	if !(reaction == 1 || reaction == 0) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if b {
		if e == reaction {
			c.deleteCommReaction(userID, commID)
		} else {
			c.updateCommReaction(userID, commID, reaction)
		}
	} else {
		c.writeCommReaction(userID, commID, reaction)
	}
	postID := c.readPostID(commID)
	url := "/post/" + strconv.Itoa(postID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (c *AppContext) readPostID(commID int) int {
	var postID int
	row := c.db.QueryRow(`SELECT post_id FROM comments WHERE comment_id = ?;`, commID)
	err := row.Scan(&postID)
	CheckErr(err)
	return postID

}

func (c *AppContext) writeCommReaction(userID, commID, reaction int) {
	_, err := c.db.Exec(`INSERT INTO comment_reaction(user_id, comment_id, reaction) VALUES(?, ?, ?)`, userID, commID, reaction)
	CheckErr(err)
}
func (c *AppContext) updateCommReaction(userID, commID, reaction int) {
	_, err := c.db.Exec(`UPDATE comment_reaction SET reaction = ? WHERE user_id = ? AND comment_id = ?;`, reaction, userID, commID)
	CheckErr(err)
}

func (c *AppContext) deleteCommReaction(userID, commID int) {
	_, err := c.db.Exec(`DELETE FROM comment_reaction WHERE user_id = ? AND comment_id = ?;`, userID, commID)
	CheckErr(err)
}

func (c *AppContext) hasReactionComm(userID, commID int) (bool, int) {
	var r int
	row := c.db.QueryRow(`SELECT reaction FROM comment_reaction WHERE user_id = ? AND comment_id = ?;`, userID, commID)
	err := row.Scan(&r)
	if err == sql.ErrNoRows {
		return false, -1
	}
	return true, r
}
