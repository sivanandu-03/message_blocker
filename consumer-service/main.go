package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)
            
func main() {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	conn, _ := amqp.Dial(os.Getenv("BROKER_URL"))
	ch, _ := conn.Channel()
	q, _ := ch.QueueDeclare("order-events", true, false, false, false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		var ev map[string]interface{}
		json.Unmarshal(d.Body, &ev)

		// Example: Update Customer LTV View
		db.Exec(`
			INSERT INTO customer_ltv_view (customer_id, total_spent, order_count, last_order_date)
			VALUES ($1, $2, 1, $3)
			ON CONFLICT (customer_id) DO UPDATE SET
			total_spent = customer_ltv_view.total_spent + EXCLUDED.total_spent,
			order_count = customer_ltv_view.order_count + 1,
			last_order_date = EXCLUDED.last_order_date
		`, ev["customerId"], 0, ev["timestamp"])

		db.Exec("UPDATE sync_status SET last_event_time = $1 WHERE id = 1", ev["timestamp"])
		log.Printf("Processed event for order %v", ev["orderId"])
	}
}