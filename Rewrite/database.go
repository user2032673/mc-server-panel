package main

import "database/sql"

func createServer(name string, maxCPU uint, maxRAM uint, maxDisk uint, maxProcs uint, db *sql.DB) (bool, uint) {
	// stub, the int is for serverID
	return false, 0
}

func deleteServer(id uint, db *sql.DB) bool {
	return false
}
