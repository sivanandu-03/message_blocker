package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type OrderItem struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	CustomerID int         `json:"customerId"`
	Items      []OrderItem `json:"items"`
}

func main() {

	connStr := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", health)

	http.HandleFunc("/api/products", createProduct)
	http.HandleFunc("/api/orders", createOrder)

	log.Println("Command Service running on 8080")

	http.ListenAndServe(":8080", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func createProduct(w http.ResponseWriter, r *http.Request) {

	var p Product

	json.NewDecoder(r.Body).Decode(&p)

	var id int

	err := db.QueryRow(`
INSERT INTO products(name, category, price, stock)
VALUES($1,$2,$3,$4)
RETURNING id`,
		p.Name, p.Category, p.Price, p.Stock).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"productId": id,
	})
}

func createOrder(w http.ResponseWriter, r *http.Request) {

	var o Order

	json.NewDecoder(r.Body).Decode(&o)

	tx, _ := db.Begin()

	var orderID int

	tx.QueryRow(`
INSERT INTO orders(customer_id,total)
VALUES($1,0)
RETURNING id`, o.CustomerID).Scan(&orderID)

	total := 0.0

	for _, item := range o.Items {

		tx.Exec(`
INSERT INTO order_items(order_id,product_id,quantity,price)
VALUES($1,$2,$3,$4)`,
			orderID, item.ProductID, item.Quantity, item.Price)

		total += item.Price * float64(item.Quantity)
	}

	tx.Exec(`
UPDATE orders SET total=$1 WHERE id=$2`,
		total, orderID)

	event := map[string]interface{}{
		"eventType":  "OrderCreated",
		"orderId":    orderID,
		"customerId": o.CustomerID,
		"total":      total,
	}

	payload, _ := json.Marshal(event)

	tx.Exec(`
INSERT INTO outbox(topic,payload)
VALUES('order-events',$1)`,
		payload)

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]int{
		"orderId": orderID,
	})
}
