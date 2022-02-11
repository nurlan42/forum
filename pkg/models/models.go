package models

type User struct {
	UserID              int
	Email, UserName string
	Password        []byte
}

type Post struct {
	PostID                    int
	UserID                    int
	Title, Content, Author    string
	TimeCreation              string
	Comments                  []Comment
	Categories                []string
	CommentNbr, Like, Dislike int
}

type Comment struct {
	CommID        int
	Author        string
	TimeCreation  string
	Content       string
	Like, Dislike int
}
