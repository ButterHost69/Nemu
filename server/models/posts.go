package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	CommentData     string
	CommentUsername string

	CreatedAt string
}

type Post struct {
	ObjectID 	string
	Username 	string
	Data     	string
	Category	string

	CreatedAt string
	Comments  []Comment
}

type BsonComment struct {
	CommentUsername 	string 				`bson:"commentUsername"`
	CommentContent     	string 				`bson:"commmentContent"` // Adjusted to match your BSON field name
	CreatedAt 			primitive.DateTime  `bson:"createdAt"`

}

type BsonPost struct {
	ObjectID  primitive.ObjectID 	`bson:"_id"`
	Username  string             	`bson:"username"`
	Data      string             	`bson:"content"`
	

	CreatedAt primitive.DateTime    `bson:"createdAt"`
	Comments  []BsonComment         `bson:"comments"`
}