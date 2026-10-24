package main

import (
	"database/sql"
	_ "embed"
)

//go:embed assets/lookupByID.sql
var lookupByIDQuery string

//go:embed assets/lookupByName.sql
var lookupByUsernameQuery string

type User struct {
	UserID       int
	Username     string
	PasswordHash string
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

func createSession() {
	
}