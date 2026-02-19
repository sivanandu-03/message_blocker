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

	connectDB()

	http.HandleFunc("/health", health)

	http.HandleFunc("/api/products", getProducts)
	http.HandleFunc("/api/orders", getOrders)

	log.Println("Query Service running 8081")

	http.ListenAndServe(":8081", nil)
}

func connectDB() {

	connStr := os.Getenv("DATABASE_URL")

	var err error

	for i := 0; i < 10; i++ {

		db, err = sql.Open("postgres", connStr)

		if err == nil {

			err = db.Ping()

			if err == nil {
				return
			}
		}

		log.Println("Waiting for DB...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal(err)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func getProducts(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.Query(`
SELECT id,name,category,price,stock
FROM products`)

	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {

		var id int
		var name, category string
		var price float64
		var stock int

		rows.Scan(&id, &name, &category, &price, &stock)

		result = append(result, map[string]interface{}{
			"id":       id,
			"name":     name,
			"category": category,
			"price":    price,
			"stock":    stock,
		})
	}

	json.NewEncoder(w).Encode(result)
}

func getOrders(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.Query(`
SELECT id,customer_id,total
FROM orders`)

	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {

		var id, customer int
		var total float64

		rows.Scan(&id, &customer, &total)

		result = append(result, map[string]interface{}{
			"id":         id,
			"customerId": customer,
			"total":      total,
		})
	}

	json.NewEncoder(w).Encode(result)
}
