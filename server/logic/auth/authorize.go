package auth

import (
	"database/sql"

	sessions "example/one-page/server/logic/session"
)

func AuthorizeUser(db *sql.DB, sessionToken string) (bool, error){
	return sessions.VerifySessionToken(db, sessionToken)
}