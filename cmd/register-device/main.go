package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func main() {
	deviceID := flag.String("id", "", "Unique Device ID (e.g. esp32-4in2-01)")
	rawKey := flag.String("key", "", "Raw API key to be hardcoded in firmware")
	dbPath := flag.String("db", ".db/app.db", "Path to SQLite database")
	flag.Parse()

	if *deviceID == "" || *rawKey == "" {
		log.Fatalf("Usage: go run cmd/register-device/main.go -id <device_id> -key <api_key>")
	}

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	keyHash := hashKey(*rawKey)

	query := `INSERT INTO devices (device_id, api_key_hash, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);`
	_, err = db.Exec(query, *deviceID, keyHash)
	if err != nil {
		log.Fatalf("Failed to insert device: %v", err)
	}

	fmt.Printf("Device successfully provisioned!\n")
	fmt.Printf("Device ID : %s\n", *deviceID)
	fmt.Printf("API Key   : %s\n", *rawKey)
}
