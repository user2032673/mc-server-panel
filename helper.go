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

	// Returns true if it exists and is a regular file (not a directory)
	return !info.IsDir()
}

func death(msg string, err error) {
	log.Fatalf("%s: %s\n", msg, err)
}
