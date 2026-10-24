package main

import (
	"database/sql"
	_ "embed"
	"uuid"
)

//go:embed assets/lookupByID.sql
var lookupByIDQuery string

//go:embed assets/lookupByName.sql
var lookupByUsernameQuery string

//go:embed assets/createSession.sql
var createSessionQuery string

//go:embed assets/validateSession.sql
var validateSessionQuery string

type User struct {
	UserID       int
	Username     string
	PasswordHash string
}

type Session struct {
	SessionID string
	UserID    int
}

func checkConnection(db *sql.DB, msg string) {
	// If the connection to the first arg DB is fine, do nothing, else fatally die with an error source of the string in arg2.
	err := db.Ping()
	if err != nil {
		cooked(msg, err)
	}
}

func userByID(id int, db *sql.DB) (bool, *User) {
	// Looks up a user given their UserID. Returns (userValid, *User). If user is invalid(likely not found),
	// then *User will be nil
	checkConnection(db, "userByID")
	row := db.QueryRow(lookupByIDQuery, id)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return false, nil
		} else {
			cooked("userByID.queryRow", row.Err())
		}
	}
	var username string
	var passwordHash string
	row.Scan(&username, &passwordHash)

	return true, &User{
		UserID:       id,
		Username:     username,
		PasswordHash: passwordHash,
	}
}

func userIDByUsername(username string, db *sql.DB) int {
	checkConnection(db, "userByName")
	row := db.QueryRow(lookupByUsernameQuery, username)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return -1
		} else {
			cooked("userByName.queryRow", row.Err())
		}
	}
	var userID int
	row.Scan(&userID)
	return userID
}

func createSession(user *User, db *sql.DB) *Session {
	checkConnection(db, "createSession")
	sessionID := uuid.NewV4().String()
	_, err := db.Exec(createSessionQuery, sessionID, user.UserID)
	if err != nil {
		cooked("createSession.createRow", err)
	}
	return &Session{
		SessionID: sessionID,
		UserID:    user.UserID,
	}

}

func validateSession(id string, db *sql.DB) (bool, *Session) {
	checkConnection(db, "validateSession")
	row := db.QueryRow(validateSessionQuery, id)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return false, nil
		} else {
			cooked("validateSession.queryRow", row.Err())
		}
	}
	var userID int
	row.Scan(userID)
	return true, &Session{
		SessionID: id,
		UserID:    userID,
	}

}
