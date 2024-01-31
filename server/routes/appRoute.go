package routes

import (
	"database/sql"
	"example/one-page/server/session"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

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