package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	connStr := os.Getenv("DATABASE_URL")

	var err error

	// retry connection until DB ready
	for i := 0; i < 10; i++ {

		db, err = sql.Open("postgres", connStr)

		if err == nil {

			err = db.Ping()

			if err == nil {
				break
			}
		}

		log.Println("Waiting for DB...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Consumer started")

	for {

		process()

		time.Sleep(3 * time.Second)
	}
}

func process() {

	rows, err := db.Query(`
SELECT id,payload FROM outbox
WHERE published_at IS NULL`)

	if err != nil {

		log.Println("Query error:", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var payload []byte

		err := rows.Scan(&id, &payload)

		if err != nil {
			continue
		}

		var event map[string]interface{}

		json.Unmarshal(payload, &event)

		log.Println("Processing event:", event)

		db.Exec(`
UPDATE outbox
SET published_at = NOW()
WHERE id = $1`, id)
	}
}
