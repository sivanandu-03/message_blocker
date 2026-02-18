package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var db *sql.DB

type Product struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type Order struct {
	CustomerID int `json:"customerId"`
	Items      []struct {
		ProductID int     `json:"productId"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	} `json:"items"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	http.HandleFunc("/api/products", createProduct)
	http.HandleFunc("/api/orders", createOrder)

	go outboxRelay() // Background poller

	log.Println("Command Service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	json.NewDecoder(r.Body).Decode(&p)
	var id int
	err := db.QueryRow("INSERT INTO products (name, category, price, stock) VALUES ($1, $2, $3, $4) RETURNING id", p.Name, p.Category, p.Price, p.Stock).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"productId": id})
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	json.NewDecoder(r.Body).Decode(&o)

	tx, _ := db.Begin()
	defer tx.Rollback()

	var orderID int
	tx.QueryRow("INSERT INTO orders (customer_id, total) VALUES ($1, $2) RETURNING id", o.CustomerID, 0).Scan(&orderID)

	// In a real app, calculate total and add items here
	payload, _ := json.Marshal(map[string]interface{}{
		"eventType": "OrderCreated",
		"orderId":   orderID,
		"customerId": o.CustomerID,
		"timestamp":  time.Now(),
	})

	tx.Exec("INSERT INTO outbox (topic, payload) VALUES ($1, $2)", "order-events", payload)
	tx.Commit()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"orderId": orderID})
}

func outboxRelay() {
	conn, _ := amqp.Dial(os.Getenv("BROKER_URL"))
	ch, _ := conn.Channel()
	for {
		rows, _ := db.Query("SELECT id, topic, payload FROM outbox WHERE published_at IS NULL LIMIT 10")
		for rows.Next() {
			var id, topic, payload string
			rows.Scan(&id, &topic, &payload)
			ch.Publish("", topic, false, false, amqp.Publishing{ContentType: "application/json", Body: []byte(payload)})
			db.Exec("UPDATE outbox SET published_at = NOW() WHERE id = $1", id)
		}
		time.Sleep(2 * time.Second)
	}
}