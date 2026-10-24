package main

import (
	"log"
	"os"
)

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // File doesn't exist
	}
	if err != nil {
		return false // Some other error (e.g., permission denied)
	}
	return !info.IsDir() // make sure it's not say a directory called app.db, which would die on .read().
}

func cooked(msg string, err error) {
	log.Fatalf("%s: %s\n", msg, err)
}
