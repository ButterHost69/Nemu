package auth

import (
	"database/sql"
	"errors"
	"fmt"

	database "example/one-page/server/db"
	sessions "example/one-page/server/logic/session"

	"github.com/go-sql-driver/mysql"
)

func CreateUser(db *sql.DB, username string, password string) error {
	
	if err := database.CreateUser(db, username, password); err != nil {
		// err.Error() == "Error 1062 (23000): Duplicate entry 'palas' for key 'users.PRIMARY'"
		if mysqlerr, ok := err.(*mysql.MySQLError); ok {
			fmt.Println(" Username already Exists Register !! ")
			if mysqlerr.Number == 1062 {
				return errors.New("User Exists")
			} else {
				return errors.New("Internal Error")
			}
		}
		return err
	} else {
		fmt.Println(" New User Registered ")
		return nil
	}
}

func LoginUser(db *sql.DB, username string, password string) (bool, string, error) {

	if database.CheckIfUserExists(db, username, password) {
		
		// Generate Tokens
		ifSessionToken, newSessionToken := sessions.CreateSessionTokens(db, username)
		if ifSessionToken == false {
			return false, "", errors.New("Internal Error")
		} else {
			return true, newSessionToken, nil
		}
		
	} else {
		return false,"", nil	
	}

}

func SignOutUser(db *sql.DB, token string){

	
	sessions.DeleteSession(db, token)
	
}

