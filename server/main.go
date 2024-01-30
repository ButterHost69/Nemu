package main

import (
	"database/sql"
	"fmt"
	"html/template"

	"net/http"

	database "example/one-page/server/db"
	routes "example/one-page/server/routes"
	sessions "example/one-page/server/session"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/static"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)


type Users struct {
	username string
	password string
}


func defaultRoute(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(c.Writer, nil)
}

func showLoginPage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("templates/loginTemplate.html"))
	tmpl.Execute(c.Writer, nil)
}

// Modify DB's
func createUser(c *gin.Context, db *sql.DB){	


	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := database.CreateUser(db, username, password); err != nil{
		// err.Error() == "Error 1062 (23000): Duplicate entry 'palas' for key 'users.PRIMARY'"
		if mysqlerr, ok := err.(*mysql.MySQLError); ok{
				fmt.Println(" Username already Exists Register !! ")
				if mysqlerr.Number == 1062{
					c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>🐦  User Exists </div>")
					} else {
						c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>😫  Internal Error </div>")
					}
			} 
		} else {
			fmt.Println(" New User Registered ")
			c.String(http.StatusOK, "<div id	='resultMessage' class='flex bg-green-200 border-2 pl-5 pr-5'>😊  Registered Successfully </div>")
		}	
}

func loginUser(c *gin.Context, db *sql.DB){

	// Check if Already Logged In...
	sessionToken, serr := c.Cookie("user-token")
	if serr == nil {
		// Check if Token is found by checking if err is null
		fmt.Println("Token Present ??") 
		if sessions.VerifySessionToken(db, sessionToken){
			fmt.Print("Redirecting to Main Page As Token Exists")
			templ := template.Must(template.ParseFiles("templates/app.html"))
			templ.Execute(c.Writer, nil)
		}
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	if database.CheckIfUserExists(db, username, password) {
		// c.String(http.StatusOK, "<div id='resultMessage' class='flex bg-green-200 border-2 pl-5 pr-5'>😁 Logged in Successfully </div>")
		
		// Generate Tokens
		ifSessionToken, newSessionToken := sessions.CreateSessionTokens(db, username)
		if ifSessionToken == false {
			c.String(http.StatusNonAuthoritativeInfo, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>😫 Internal Error Occured </div>")
		} 
		
		c.SetCookie("user-token", newSessionToken, 3600, "/", "localhost", false, true)
		// c.SetCookie("user-token",username,3600,"/","localhost",false,true)
		
		fmt.Println("Cookie Set", newSessionToken)
		
		templ := template.Must(template.ParseFiles("templates/app.html"))
		templ.Execute(c.Writer, nil)
	} else {
		c.String(http.StatusNonAuthoritativeInfo, "<div id='resultMessage' class='flex bg-red-200 border-2 pl-5 pr-5'>🤦‍♂️ Incorrect Username/Password </div>")
	}
}

func main() {
	fmt.Println("  !!! Hello World !!! ")

	r := gin.Default()
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Error > ", err.Error())
		return
	}

	//r.Use(static.Serve("/",static.LocalFile("/js/", true)))
	r.Static("/js","./js")

	r.GET("/", defaultRoute)
	r.GET("/login", showLoginPage)
	r.GET("/app", routes.AppHomeRoute)
	

	r.POST("/", func(c *gin.Context) {
		createUser(c, db)
	})
	r.POST("/login", func(c *gin.Context) {
		loginUser(c, db)
	})

	r.Run("localhost:8000")
}