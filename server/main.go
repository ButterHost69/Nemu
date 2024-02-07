package main

import (
	"fmt"
	"log"

	database "example/one-page/server/db"
	handler "example/one-page/server/handler"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/static"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	username string
	password string
}


func main() {
	fmt.Println("  !!! Hello World !!! ")

	r := gin.Default()

	// Init MYSQL Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error > ", err.Error())
		return
	}

	// Init MongoDB Database
	mdb, merr := database.InitMongoDB()
	if merr != nil {
		log.Fatal("Error > ", err.Error())
		return
	}

	// Root Page
	r.GET("/", func(ctx *gin.Context) {
		handler.DefaultRoute(ctx)
	})

	r.GET("/app", func(ctx *gin.Context) {
		handler.GetAppPage(db, ctx)
	})

	r.GET("/login", func(ctx *gin.Context) {
		handler.GetLoginPage(ctx)
	})

	r.GET("/register", func(ctx *gin.Context) {
		handler.GetRegisterPage(ctx)
	})



	// Creates User
	r.POST("/register", func(c *gin.Context) {
		handler.PostRegister(db, c)
	})

	r.POST("/login", func(c *gin.Context) {
		handler.PostLogin(c, db)
	})

	r.POST("/signout", func(c *gin.Context) {
		handler.PostSignOut(c, db)
	})

	// Create Posts
	r.POST("/app/post", func(ctx *gin.Context) {
		handler.PostAppPost(ctx, mdb, db)
	})

	r.POST("/app/comment/:postId", func(ctx *gin.Context) {
		handler.PostComment(ctx, mdb, db)
	})

	r.GET("/app/posts/:page", func(ctx *gin.Context) {
		handler.LoadPages(ctx, mdb)
	})

	r.GET("/app/components/commentBox/:postID", func(ctx *gin.Context) {
		// ctx.String(200,"Pressed")
		handler.GetCommentInputBox(ctx, mdb)
	})

	r.GET("/app/categories/:category", func(ctx *gin.Context) {
		handler.GetCategoryPage(ctx)
	})

	r.POST("/app/:category/post", func(ctx *gin.Context) {
		handler.PostCreateCategoryPost(ctx, mdb, db)
	})

	r.Run("localhost:8000")
}
