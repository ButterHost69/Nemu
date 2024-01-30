package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	database "example/one-page/server/db"
	"fmt"
)

func VerifySessionToken(db *sql.DB, token string)(bool){
	return checkIfSessionExists(db, token)
}

func CreateSessionTokens(db *sql.DB, userName string)(bool, string){
	// If checkIfUserExist return true -> than user exists in the db
	// Returns True If Created
	// Returns False if Not Created

	if checkIfUserExists(db, userName){
		return false, ""
	}

	token, err := generateToken()

	if err != nil{
		fmt.Println(" > Error In Generate Token : ", err.Error())
		return false, ""
	}

	// If checkIfSessionExists return true -> than Session Id is a Duplicate
	if checkIfSessionExists(db, token){
		fmt.Println(" > Error In Check If Session Exists Token : ", err.Error())
		return false, ""
	}

	if err = storeInDB(db, userName, token) ; err != nil {
		fmt.Println(" > Error In Store In Db Token : ", err.Error())
		return	false, ""
	}
	
	return true, token
}

func generateToken()(string, error){
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	customEncoding := base64.RawURLEncoding
	token := customEncoding.EncodeToString(randomBytes)

	return token, nil
}

func storeInDB(db *sql.DB, username string, sessionToken string)error{
	if err := database.CreateSession(db, username, sessionToken); err != nil {
		return err
	}

	return nil
}

func checkIfUserExists(db *sql.DB, username string)(bool){
	return database.CheckIfUserExistsUsingUsernameInSessionDB(db, username)
}

func checkIfSessionExists(db *sql.DB, sessionId string)(bool){
	return database.CheckIfSessionIdExistsUsingSessionIdInSessionDB(db, sessionId)
}

