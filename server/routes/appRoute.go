package routes

import (
	
	"database/sql"
	"example/one-page/server/session"
	"fmt"
	"html/template"
	

	"github.com/gin-gonic/gin"
	
	"go.mongodb.org/mongo-driver/mongo"


	database "example/one-page/server/db"
)

type Comment struct {
	CommentData string
	CommentUsername string
}

type Post struct{
	Username string
	Data string
	
	Comments []Comment
}

func AppHomeRoute(c *gin.Context){
		
}

func AppSignout(c *gin.Context, db *sql.DB){
	token, err := c.Cookie("user-token")
	if err != nil {
		fmt.Println(" > Error Occured At SignOut : " , err.Error())
		return
	}
	session.DeleteSession(db, token)
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(c.Writer, nil)
}

func CreatePost(c *gin.Context, mdb *mongo.Client, db *sql.DB){
	postCollection := mdb.Database("nemu").Collection("posts")
	
	user_token, err := c.Cookie("user-token")
	if err!=nil {
		fmt.Println("Error Occured in Getting Session-Token [Create Post] : ")
		fmt.Println(err.Error())
		c.String(200, "<p> Unable to Create Post </p>")

		return	
	}

	username, derr := session.GetUsernameFromSessionToken(db, user_token)
	if derr != nil {
		fmt.Println("Error Occured in Getting Username From Session [Create Post] : ")
		fmt.Println(derr.Error())
		c.String(200, "<p> Unable to Create Post </p>")

		return	
	}

	content := c.PostForm("post-content")
	if content == ""{
		fmt.Println("Could Not Retrieve any Content From the Form [Create Post] ")
		c.String(200, "<p> Unable to Create Post </p>")

		return	
	}

	merr := database.InsertPostInMongoDB(postCollection, username, content)
	if merr != nil{
		fmt.Println("Error Occured in Inserting Post in DB [Create Post] : ")
		fmt.Println(merr.Error())
		c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	currPost := Post {
		Username: username,
		Data: content,

		// Temp Delete Later :
		// Comments: []Comment{
		// 	{
		// 		CommentData : "Comment - 1",
		// 		CommentUsername : "Well Well",
		// 	},
		// },
	}

	tmplt, terr := template.ParseFiles("templates/posts.html")
	if terr != nil {
		fmt.Println("Error Occured in Parse HTML template [Create Post] : ")
		fmt.Println(terr.Error())
		c.String(200, "<p> Unable to Display Post </p>")
		return
	}

	tmplt.Execute(c.Writer, currPost)
	// c.String(200, "<p> New Post Created </p>")
}