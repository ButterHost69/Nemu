package posts

import (
	"database/sql"
	"fmt"
	
	
	// "strconv"
	
	"go.mongodb.org/mongo-driver/mongo"
	
	database "example/one-page/server/db"
	models "example/one-page/server/models"
	sessions "example/one-page/server/logic/session"
)


func CreatePost(mdb *mongo.Client, db *sql.DB, user_token string, content string) (models.Post, error){
	postCollection := mdb.Database("nemu").Collection("posts")

	
	emptyPost := models.Post{
		Username: "",
		Data: "",
	}

	
	username, derr := sessions.GetUsernameFromSessionToken(db, user_token)
	if derr != nil {
		fmt.Println("Error Occured in Getting Username From Session [Create Post] : ")
		fmt.Println(derr.Error())
		
		return emptyPost, derr
	}

	merr := database.InsertPostInMongoDB(postCollection, username, content)
	if merr != nil {
		fmt.Println("Error Occured in Inserting Post in DB [Create Post] : ")
		fmt.Println(merr.Error())
		
		return emptyPost, merr
	}

	currPost := models.Post{
		Username: username,
		Data:     content,

		// Temp Delete Later :
		// Comments: []Comment{
		// 	{
		// 		CommentData : "Comment - 1",
		// 		CommentUsername : "Well Well",
		// 	},
		// },
	}

	return currPost, nil
}

func GetPosts(mdb *mongo.Client, pageNumber int) ([]models.Post, bool){

	postsCollection := mdb.Database("nemu").Collection("posts")
	

	// This func Returns models.Post[]
	return database.GetPostsFromMongoDB(postsCollection, pageNumber*10, 10)

}