package main

import (
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"strings"

	incus "github.com/lxc/incus/client"
	_ "modernc.org/sqlite"
)

//go:embed assets/firstRun.sql
var firstRunSQL string

//go:embed assets/bindAddress.txt
var bindAddress string

func firstRunSetup() *sql.DB {
	log.Println("First run setup . . .")
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		cooked("dbOpen", err)
	}
	_, err = db.Exec(firstRunSQL)
	if err != nil {
		cooked("dbSetup", err)
	}
	return db
}

func main() {
	log.Println("Minecraft Server Panel v1.0 starting . . .")
	client, err := incus.ConnectIncusUnix("", nil)
	if err != nil {
		cooked("connect", err)
	}
	serverInfo, _, err := client.GetServer()
	if err != nil {
		cooked("getServer", err)
	}
	if serverInfo.APIVersion != "1.0" {
		log.Fatalln("Wrong Incus server API version, refusing to connect for now!")
	}
	if !isFile("app.db") {
		db := firstRunSetup()
		firstRunSetup()
	} else {
		db, err := sql.Open("sqlite", "app.db")
		if err != nil {
			cooked("dbOpen", err)
		}
	}
	// http is a global object here so the function doesn't lose track of it
	ramDB := loadRAMDB(db)
	httpMethods(ramDB, db) // should be called setupHttpMethods
	log.Printf("Binding to %s . . . \n", strings.TrimSpace(bindAddress))
	err = http.ListenAndServe(strings.TrimSpace(bindAddress), nil)
	if err != nil {
		cooked("listenAndServe", err)
	}
}
