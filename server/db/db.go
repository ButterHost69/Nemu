package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	dblink := "root:deep1520@tcp(127.0.0.1:3306)/nemu"
	db, err := sql.Open("mysql", dblink)
	if err != nil {
		return nil, err
	}

	// Pinging The Database To Verify The Connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(db *sql.DB, username string, password string) error {
	query := "INSERT INTO users VALUES (?,?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errs := stmt.Exec(username, password)
	if errs != nil {
		return errs
	}

	return nil
}

func CheckIfUserExistsUsingUsernameInSessionDB(db *sql.DB, username string) bool {
	query := "SELECT * FROM session_token_table WHERE username=?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return false
	}
	defer stmt.Close()

	errs := stmt.QueryRow(username).Scan(&username)

	if errs != nil {
		return false
	} else {
		return true
	}
}

func CheckIfSessionIdExistsUsingSessionIdInSessionDB(db *sql.DB, sessionToken string) bool {
	query := "SELECT * FROM session_token_table WHERE session_token=?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return false
	}
	defer stmt.Close()

	errs := stmt.QueryRow(sessionToken).Scan(&sessionToken)

	if errs != nil {
		return false
	} else {
		return true
	}
}

func CheckIfUserExists(db *sql.DB, username string, password string) bool {
	query := "SELECT * FROM users WHERE username=? AND password=?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return false
	}
	defer stmt.Close()

	errs := stmt.QueryRow(username, password).Scan(&username, &password)

	if errs != nil {
		return false
	} else {
		return true
	}
}

func CreateSession(db *sql.DB, username string, token string) error {
	query := "INSERT INTO session_token_table VALUES (?,?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errs := stmt.Exec(username, token)
	if errs != nil {
		return errs
	}

	return nil
}
