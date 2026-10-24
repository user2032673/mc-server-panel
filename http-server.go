package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
)

type LoginRequest struct {
	username string
	password string
}

func methodLock(methodType string, writer http.ResponseWriter, request *http.Request) bool {
	if request.Method != methodType {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func createSession(userID int, db *sql.DB) {

}

func httpMethods(database *sql.DB) {

	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		methodLock("GET", w, r)
		fmt.Fprint(w, "1.0.0_alpha")
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		methodLock("POST", w, r)
		var loginReq LoginRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&loginReq)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		userID := userIDByUsername(loginReq.username, database)
		if userID <= 0 {
			http.Error(w, `{"success": false, "reason": "Username or password incorrect"}`, http.StatusUnauthorized)
			return
		}
		valid, user := userByID(userID, database)
		if !valid {
			http.Error(w, `{"success": false, "reason": "User disappeared"}`, http.StatusInternalServerError)
			return
		}
		match, err := argon2id.ComparePasswordAndHash(loginReq.password, user.PasswordHash)
		if err != nil {
			cooked("verifyHash", err)
		}
		if match {
			// Happy path :)

		} else {
			// Bad path :O
		}

	})

}
