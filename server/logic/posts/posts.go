package posts

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	// "strconv"

	"go.mongodb.org/mongo-driver/mongo"

	database "example/one-page/server/db"
	sessions "example/one-page/server/logic/session"
	models "example/one-page/server/models"
)


func CreatePost(mdb *mongo.Client, db *sql.DB, user_token string, content string) (models.Post, error){
	postCollection := mdb.Database("nemu").Collection("posts")

	emptyPost := models.Post{
		Username: "",
		Data: "",
		CreatedAt: "",
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

	year, month, day := time.Now().Date()
	date := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
	currPost := models.Post{
		Username: username,
		Data:     content,
		CreatedAt: date,
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

func CreateCategoryPost(mdb *mongo.Client, db *sql.DB, user_token string, content string, category string) (models.Post, error){
	postCollection := mdb.Database("nemu").Collection("posts")

	emptyPost := models.Post{
		Username: "",
		Data: "",
		CreatedAt: "",
	}

	
	username, derr := sessions.GetUsernameFromSessionToken(db, user_token)
	if derr != nil {
		fmt.Println("Error Occured in Getting Username From Session [Create Post] : ")
		fmt.Println(derr.Error())
		
		return emptyPost, derr
	}

	merr := database.InsertCategoryPostInMongoDB(postCollection, username, content, category)
	if merr != nil {
		fmt.Println("Error Occured in Inserting Post in DB [Create Post] : ")
		fmt.Println(merr.Error())
		
		return emptyPost, merr
	}

	year, month, day := time.Now().Date()
	date := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
	currPost := models.Post{
		Username: username,
		Data:     content,
		CreatedAt: date,
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

func CreateComment(mdb *mongo.Client, db *sql.DB, user_token string, postObjId string, commentContent string) (models.Comment, error) {
	postCollection := mdb.Database("nemu").Collection("posts")

	emptyPost := models.Comment{
		CommentUsername: "",
		CommentData: "",
		CreatedAt: "",
	}

	
	username, derr := sessions.GetUsernameFromSessionToken(db, user_token)
	if derr != nil {
		fmt.Println("Error Occured in Getting Username From Session [Insert Comment] : ")
		fmt.Println(derr.Error())
		
		return emptyPost, derr
	}

	merr := database.InsertCommentInMongoDB(postCollection, postObjId, username, commentContent)
	if merr != nil {
		fmt.Println("Error Occured in Inserting Post in DB [Create Post] : ")
		fmt.Println(merr.Error())
		
		return emptyPost, merr
	}

	year, month, day := time.Now().Date()
	date := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
	currPost := models.Comment{
		CommentData: commentContent,
		CommentUsername: username,
		CreatedAt: date,
	}

	return currPost, nil
}

func GetPosts(mdb *mongo.Client, pageNumber int) ([]models.Post, bool){

	postsCollection := mdb.Database("nemu").Collection("posts")
	

	// This func Returns models.Post[]
	return database.BetterGetPostsFromMongoDB(postsCollection, pageNumber*10, 10)

}

func GetCategoryPosts (mdb *mongo.Client, pageNumber int, category string) ([]models.Post, bool){

	postsCollection := mdb.Database("nemu").Collection("posts")
	

	// This func Returns models.Post[]
	return database.BetterGetCategoryPostsFromMongoDB(postsCollection, pageNumber*10, 10, category)

}

func GetUsernameThroughObjectID(mdb *mongo.Client, objId string) (string) {
	postsCollection := mdb.Database("nemu").Collection("posts")

	return database.GetUsernameThroughObjectID(postsCollection, objId)
}