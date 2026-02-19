package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	_ "github.com/lib/pq"
)

var db *sql.DB
var channel *amqp.Channel

func main() {

	connectDB()
	connectRabbit()

	log.Println("Consumer started")

	for {

		process()

		time.Sleep(2 * time.Second)
	}
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

func connectRabbit() {

	url := "amqp://guest:guest@broker:5672/"

	var err error

	for i := 0; i < 10; i++ {

		conn, err := amqp.Dial(url)

		if err == nil {

			channel, err = conn.Channel()

			if err == nil {
				return
			}
		}

		log.Println("Waiting for RabbitMQ...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal(err)
}

func process() {

	rows, err := db.Query(`
SELECT id, topic, payload
FROM outbox
WHERE published_at IS NULL`)

	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var topic string
		var payload []byte

		rows.Scan(&id, &topic, &payload)

		err := channel.Publish(
			"",
			topic,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        payload,
			},
		)

		if err != nil {
			continue
		}

		db.Exec(`
UPDATE outbox
SET published_at = NOW()
WHERE id=$1`, id)

		log.Println("Published event:", id)
	}
}
