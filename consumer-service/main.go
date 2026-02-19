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

var db *sql.DB

// Event structure from outbox
type OrderCreatedEvent struct {
	EventType  string  `json:"eventType"`
	OrderID    int     `json:"orderId"`
	CustomerID int     `json:"customerId"`
	Total      float64 `json:"total"`
}

func connectDB() *sql.DB {

	connStr := os.Getenv("DATABASE_URL")

	for {
		db, err := sql.Open("postgres", connStr)

		if err != nil {
			log.Println("Waiting for DB...")
			time.Sleep(3 * time.Second)
			continue
		}

		err = db.Ping()

		if err != nil {
			log.Println("Waiting for DB...")
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("Connected to DB")
		return db
	}
}

func connectRabbitMQ() *amqp.Connection {

	url := "amqp://guest:guest@broker:5672/"

	for {
		conn, err := amqp.Dial(url)

		if err != nil {
			log.Println("Waiting for RabbitMQ...")
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("Connected to RabbitMQ")
		return conn
	}
}

func main() {

	// Connect DB
	db = connectDB()
	defer db.Close()

	// Connect RabbitMQ
	conn := connectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal("Channel error:", err)
	}

	defer ch.Close()

	// Declare Queue
	q, err := ch.QueueDeclare(
		"order-events",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Consume error:", err)
	}

	log.Println("Consumer Service Started")

	forever := make(chan bool)

	go func() {

		for msg := range msgs {

			var event OrderCreatedEvent

			err := json.Unmarshal(msg.Body, &event)

			if err != nil {
				log.Println("JSON error:", err)
				continue
			}

			log.Println("Received Order:", event.OrderID)

			// Insert into read model
			_, err = db.Exec(`
				INSERT INTO orders_read(order_id, customer_id, total)
				VALUES ($1,$2,$3)
				ON CONFLICT (order_id) DO NOTHING
			`,
				event.OrderID,
				event.CustomerID,
				event.Total,
			)

			if err != nil {
				log.Println("DB insert error:", err)
				continue
			}

			log.Println("Saved to read DB:", event.OrderID)
		}

	}()

	<-forever
}
