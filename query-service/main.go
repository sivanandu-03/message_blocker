package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	conn := os.Getenv("DATABASE_URL")

	db, _ = sql.Open("postgres", conn)

	http.HandleFunc("/health", health)

	http.HandleFunc("/api/analytics/products/", productSales)

	log.Println("Query Service running 8081")

	http.ListenAndServe(":8081", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func productSales(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/api/analytics/products/"):]
	id, _ := strconv.Atoi(idStr)

	var qty int
	var revenue float64
	var count int

	db.QueryRow(`
SELECT total_quantity_sold,total_revenue,order_count
FROM product_sales_view
WHERE product_id=$1`,
		id).Scan(&qty, &revenue, &count)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"productId":         id,
		"totalQuantitySold": qty,
		"totalRevenue":      revenue,
		"orderCount":        count,
	})
}
