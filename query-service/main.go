package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	http.HandleFunc("/api/analytics/sync-status", getSyncStatus)
	// Add other analytics handlers here...

	log.Println("Query Service starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getSyncStatus(w http.ResponseWriter, r *http.Request) {
	var lastTime time.Time
	db.QueryRow("SELECT last_event_time FROM sync_status WHERE id = 1").Scan(&lastTime)
	
	lag := time.Since(lastTime).Seconds()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"lastProcessedEventTimestamp": lastTime.Format(time.RFC3339),
		"lagSeconds":                lag,
	})
}