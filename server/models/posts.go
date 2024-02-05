package models


type Comment struct {
	CommentData     string
	CommentUsername string

	CreatedAt 		string
}

type Post struct {
	ObjectID 	string
	Username 	string
	Data     	string

	CreatedAt 	string
	Comments   []Comment
}