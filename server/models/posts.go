package models


type Comment struct {
	CommentData     string
	CommentUsername string
}

type Post struct {
	Username 	string
	Data     	string

	CreatedAt 	string
	Comments   []Comment
}