package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	auth "example/one-page/server/logic/auth"
	"example/one-page/server/logic/posts"
	"example/one-page/server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var domain string = "localhost"

func DefaultRoute(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("pages/index.html"))
	tmpl.Execute(c.Writer, nil)

}

func GetAppPage(db *sql.DB, c *gin.Context) {
	// Get Request to /app Page :
	// If correct Session Id Present returns -> main[functional] page
	// Else returns -> register page

	sessionToken, serr := c.Cookie("user-token")
	if serr == nil {
		// Check if Token is found by checking if err is null
		fmt.Println(" > Cookies Detected !!")
		if exists, err := auth.AuthorizeUser(db, sessionToken); err == nil && exists {
			fmt.Println(" > Redirecting to Main Page As Token Exists")
			templ := template.Must(template.ParseFiles("pages/app.html"))

			data := struct {
				Category string
			}{
				Category: "all",
			}
			err := templ.Execute(c.Writer, data)
			if err != nil {
				fmt.Println(" > Error in Executing Template : ", err.Error())
			}
			return
		}
	}

	GetRegisterPage(c)
}

func GetLoginPage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("pages/login.html"))
	tmpl.Execute(c.Writer, nil)
}

func GetRegisterPage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("pages/register.html"))
	tmpl.Execute(c.Writer, nil)
}

func PostRegister(db *sql.DB, c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := auth.CreateUser(db, username, password); err != nil {
		if err.Error() == "User Exists" {
			c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>🐦  User Exists </div>")
			return
		} else {
			c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>😫  Internal Error </div>")
			return
		}
	} else {
		c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-green-200 border-2 pl-5 pr-5'>😊  Registered Successfully </div>")
		return
	}
}

func PostLogin(c *gin.Context, db *sql.DB) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	ifAuth, newSessionToken, err := auth.LoginUser(db, username, password)
	if err != nil {
		c.String(http.StatusNonAuthoritativeInfo, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>😫 Internal Error Occured </div>")
		return
	}

	if ifAuth == true {
		c.SetCookie("user-token", newSessionToken, 3600, "/", domain, false, true)
		// c.SetCookie("user-token",username,3600,"/","localhost",false,true)

		fmt.Println("Cookie Set", newSessionToken)

		templ := template.Must(template.ParseFiles("pages/app.html"))
		templ.Execute(c.Writer, nil)
		return
	} else {
		c.String(http.StatusNonAuthoritativeInfo, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>🤦‍♂️ Incorrect Username/Password </div>")
		return
	}
}

func PostSignOut(c *gin.Context, db *sql.DB) {
	token, err := c.Cookie("user-token")
	if err != nil {
		fmt.Println(" > Error Occured At SignOut : ", err.Error())
		return
	}

	auth.SignOutUser(db, token)

	c.SetCookie("user-token", "", -1, "/", "localhost", false, true) //Deletes Cookie user-token
	tmpl := template.Must(template.ParseFiles("pages/index.html"))
	tmpl.Execute(c.Writer, nil)
}

func PostAppPost(c *gin.Context, mdb *mongo.Client, db *sql.DB) {

	user_token, err := c.Cookie("user-token")
	if err != nil {
		fmt.Println("Error Occured in Getting Session-Token [Create Post] : ")
		fmt.Println(err.Error())

		c.String(200, "<p> Unable to Create Post </p>")
	}

	content := c.PostForm("post-content")
	if content == "" {
		fmt.Println("Could Not Retrieve any Content From the Form [Create Post] ")
		// c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	currPost, err := posts.CreatePost(mdb, db, user_token, content)
	if err != nil {
		c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	tmplt, terr := template.ParseFiles("components/comments.html", "components/posts.html")
	if terr != nil {
		fmt.Println("Error Occured in Parse HTML template [Create Post] : ")
		fmt.Println(terr.Error())
		c.String(200, "<p> Unable to Display Post </p>")
		return
	}

	tmplt.ExecuteTemplate(c.Writer, "postsComponent", currPost)
	return
}

func LoadPages(c *gin.Context, mdb *mongo.Client) {
	strpageNumber := c.Param("page")

	pageNumber, err := strconv.Atoi(strpageNumber)
	if err != nil {
		// Handle the error if the conversion fails
		c.String(200, "Invalid ID")
		return
	}

	posts, ifEmptyPost := posts.GetPosts(mdb, pageNumber)

	type PageData struct {
		PageNumber int
		PostsSet   []models.Post
		EmptyPost  bool
		Category   string
	}

	tmpl, tempErr := template.ParseFiles("components/loadPosts.html", "components/posts.html", "components/comments.html")
	if tempErr != nil {
		fmt.Println("Error executing Parsing Templates : ", err.Error())
		// return
	}

	data := PageData{
		PageNumber: pageNumber + 1,
		PostsSet:   posts,
		EmptyPost:  ifEmptyPost,
	}
	if err := tmpl.ExecuteTemplate(c.Writer, "loadPosts", data); err != nil {
		// Handle the error, maybe log it or return an error response.
		fmt.Println("Error executing template:", err.Error())
		return
	}
}

func GetCommentInputBox(c *gin.Context, mdb *mongo.Client) {
	objectID := c.Param("postID")
	usernameOfOP := posts.GetUsernameThroughObjectID(mdb, objectID)

	fmt.Println(objectID)
	tmpl, err := template.ParseFiles("components/inputBox.html")
	if err != nil {
		fmt.Println(" > Error in Rendering Comment Input Box ")
		fmt.Println(err.Error())
	}

	data := struct {
		ObjectID   string
		OPUsername string
	}{
		ObjectID:   objectID,
		OPUsername: usernameOfOP,
	}
	errs := tmpl.ExecuteTemplate(c.Writer, "inputCommentBox", data)
	if errs != nil {
		fmt.Println(" > Error in Rendering Comment Input Box ")
		fmt.Println(errs.Error())
	}

	return
}

func PostComment(c *gin.Context, mdb *mongo.Client, db *sql.DB) {

	user_token, err := c.Cookie("user-token")
	if err != nil {
		fmt.Println("Error Occured in Getting Session-Token [Create Post] : ")
		fmt.Println(err.Error())

		c.String(200, "<p> Unable to Create Post </p>")
	}

	objectID := c.Param("postId")
	commentContent := c.PostForm("post-content")
	if commentContent == "" {
		fmt.Println(" > No Content in Comment Post ")
	}

	commentData, err := posts.CreateComment(mdb, db, user_token, objectID, commentContent)

	if err != nil {
		c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	tmplt, terr := template.ParseFiles("components/comments.html", "components/posts.html")
	if terr != nil {
		fmt.Println("Error Occured in Parse HTML template [Create Post] : ")
		fmt.Println(terr.Error())
		c.String(200, "<p> Unable to Display Post </p>")
		return
	}

	tmplt.ExecuteTemplate(c.Writer, "commentsComponent", commentData)
	return
}

func GetCategoryPage(c *gin.Context) {
	category := c.Param("category")
	// fmt.Println("> Category :", category)
	if category == "all" {
		templ := template.Must(template.ParseFiles("pages/app.html"))
		data := struct {
			Category string
		}{
			Category: "all",
		}
		err := templ.Execute(c.Writer, data)
		if err != nil {
			fmt.Println(" > Error in Executing Template : ", err.Error())
		}
		return
	}

	tmpl, err := template.ParseFiles("pages/categoryChange.html")
	if err != nil {
		fmt.Println(" > Error in Rendering Category Change ")
		fmt.Println(err.Error())
	}

	data := struct {
		Category string
	}{
		Category: category,
	}
	errs := tmpl.ExecuteTemplate(c.Writer, "categoryChange", data)
	if errs != nil {
		fmt.Println(" > Error in Rendering Comment Input Box ")
		fmt.Println(errs.Error())
	}
}

func PostCreateCategoryPost(c *gin.Context, mdb *mongo.Client, db *sql.DB) {
	category := c.Param("category")
	fmt.Println(" > Catgory : ", category)

	user_token, err := c.Cookie("user-token")
	if err != nil {
		fmt.Println("Error Occured in Getting Session-Token [Create Post] : ")
		fmt.Println(err.Error())

		c.String(200, "<p> Unable to Create Post </p>")
	}

	content := c.PostForm("post-content")
	if content == "" {
		fmt.Println("Could Not Retrieve any Content From the Form [Create Post] ")
		// c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	currPost, err := posts.CreateCategoryPost(mdb, db, user_token, content, category)
	if err != nil {
		c.String(200, "<p> Unable to Create Post </p>")
		return
	}

	tmplt, terr := template.ParseFiles("components/comments.html", "components/posts.html")
	if terr != nil {
		fmt.Println("Error Occured in Parse HTML template [Create Post] : ")
		fmt.Println(terr.Error())
		c.String(200, "<p> Unable to Display Post </p>")
		return
	}

	tmplt.ExecuteTemplate(c.Writer, "postsComponent", currPost)
	return
}

func LoadCategoryPages(c *gin.Context, mdb *mongo.Client) {
	strpageNumber := c.Param("page")
	category := c.Param("category")

	pageNumber, err := strconv.Atoi(strpageNumber)
	if err != nil {
		// Handle the error if the conversion fails
		c.String(200, "Invalid ID")
		return
	}

	posts, ifEmptyPost := posts.GetCategoryPosts(mdb, pageNumber, category)

	type PageData struct {
		PageNumber int
		PostsSet   []models.Post
		EmptyPost  bool
		Category   string
	}

	tmpl, tempErr := template.ParseFiles("components/categoryLoadPost.html", "components/posts.html", "components/comments.html")
	if tempErr != nil {
		fmt.Println("Error executing Parsing Templates : ", err.Error())
		// return
	}

	data := PageData{
		PageNumber: pageNumber + 1,
		PostsSet:   posts,
		EmptyPost:  ifEmptyPost,
	}
	if err := tmpl.ExecuteTemplate(c.Writer, "categoryLoadPosts", data); err != nil {
		// Handle the error, maybe log it or return an error response.
		fmt.Println("Error executing template:", err.Error())
		return
	}
}
