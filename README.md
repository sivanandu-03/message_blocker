# High-Performance E-Commerce Analytics System (CQRS + Event-Driven Architecture)

## Overview

This project implements a **high-performance e-commerce analytics backend** using:

* CQRS (Command Query Responsibility Segregation)
* Event-Driven Architecture (EDA)
* RabbitMQ Message Broker
* PostgreSQL Write and Read Models
* Docker & Docker Compose for full containerization

The system separates **write operations (commands)** from **read operations (queries)** and uses asynchronous event processing to maintain optimized read models.

This architecture ensures:

* High scalability
* High performance analytics
* Loose coupling between services
* Fault tolerance and reliability

---

## Architecture Diagram

Components:

```
Client
  │
  ▼
Command Service (Write API)
  │
  ▼
PostgreSQL (Write DB)
  │
  ▼
Outbox Table
  │
  ▼
RabbitMQ Message Broker
  │
  ▼
Consumer Service
  │
  ▼
Read Database (Materialized Views)
  │
  ▼
Query Service (Read API)
  │
  ▼
Client
```

---

## Technologies Used

| Technology     | Purpose                 |
| -------------- | ----------------------- |
| Go             | Backend Services        |
| PostgreSQL     | Write and Read Database |
| RabbitMQ       | Message Broker          |
| Docker         | Containerization        |
| Docker Compose | Service orchestration   |
| REST API       | Communication           |

---

## CQRS Pattern Implementation

### Command Side

Handles write operations:

* Create Product
* Create Order
* Writes to database
* Writes event to Outbox table

Service: `command-service`

---

### Query Side

Handles read operations:

* Returns analytics data
* Uses optimized read models
* No write operations

Service: `query-service`

---

### Event-Driven Consumer

Processes events asynchronously:

* Reads events from RabbitMQ
* Updates read models
* Ensures eventual consistency

Service: `consumer-service`

---

## Folder Structure

```
message_blocker/
│
├── docker-compose.yml
├── init.sql
├── README.md
├── submission.json
├── .env.example
│
├── command-service/
│   ├── Dockerfile
│   └── main.go
│
├── query-service/
│   ├── Dockerfile
│   └── main.go
│
├── consumer-service/
│   ├── Dockerfile
│   └── main.go
```

---

## Database Schema

### Write Model Tables

products

```
id
name
category
price
stock
```

orders

```
id
customer_id
total
created_at
```

order_items

```
id
order_id
product_id
quantity
price
```

outbox

```
id
topic
payload
created_at
published_at
```

---

### Read Model Tables

orders_read

```
order_id
customer_id
total
```

---

## Docker Setup

This project uses Docker Compose to start all services.

Services included:

* PostgreSQL
* RabbitMQ
* Command Service
* Query Service
* Consumer Service

---

## How to Run the Project

### Step 1 — Install Requirements

Install:

* Docker Desktop
* Go (optional)
* Postman (optional)

---

### Step 2 — Start System

Run:

```
docker-compose up --build
```

Expected output:

```
Command Service running on 8080
Query Service running 8081
Consumer Service Started
Connected to RabbitMQ
Connected to DB
```

---

## Service URLs

Command Service:

```
http://localhost:8080
```

Query Service:

```
http://localhost:8081
```

RabbitMQ Dashboard:

```
http://localhost:15672
```

Login:

```
username: guest
password: guest
```

---

## API Endpoints

---

### Create Product

POST

```
/api/products
```

Example Request:

```
http://localhost:8080/api/products
```

Body:

```json
{
  "name": "Laptop",
  "category": "Electronics",
  "price": 50000,
  "stock": 10
}
```

Response:

```json
{
  "productId": 1
}
```

---

### Create Order

POST

```
/api/orders
```

Example Request:

```
http://localhost:8080/api/orders
```

Body:

```json
{
  "customerId": 101,
  "items": [
    {
      "productId": 1,
      "quantity": 2,
      "price": 50000
    }
  ]
}
```

Response:

```json
{
  "orderId": 1
}
```

---

### Query Products

GET

```
http://localhost:8081/api/products
```

---

## Event Flow Explanation

Step 1
Command Service creates order

Step 2
Order saved in write database

Step 3
Event saved in outbox table

Step 4
Event published to RabbitMQ

Step 5
Consumer receives event

Step 6
Consumer updates read model

Step 7
Query Service returns analytics

---

## Event Example

OrderCreated Event:

```json
{
  "eventType": "OrderCreated",
  "orderId": 1,
  "customerId": 101,
  "total": 100000
}
```

---

## Testing Guide

### Test 1 — Verify Docker Containers

Run:

```
docker ps
```

Verify all running:

```
command-service
query-service
consumer-service
db
broker
```

---

### Test 2 — Create Product

Use Postman:

POST

```
http://localhost:8080/api/products
```

---

### Test 3 — Create Order

POST

```
http://localhost:8080/api/orders
```

---

### Test 4 — Verify Consumer Processing

Check logs:

```
docker-compose logs consumer-service
```

Expected:

```
Received Order: 1
Saved to read DB: 1
```

---

### Test 5 — Verify RabbitMQ

Open:

```
http://localhost:15672
```

---

### Test 6 — Verify Read Model

Query using Query Service

```
http://localhost:8081/api/products
```

---

## Outbox Pattern Implementation

Ensures reliable event publishing.

Process:

1. Write business data
2. Write event to outbox table
3. Consumer publishes event
4. Event processed safely

Prevents data loss.

---

## Health Check Endpoints

Command Service:

```
/health
```

Query Service:

```
/health
```

---

## submission.json

```
{
  "commandServiceUrl": "http://localhost:8080",
  "queryServiceUrl": "http://localhost:8081"
}
```

---

## .env.example

```
DATABASE_URL=postgres://user:password@db:5432/write_db?sslmode=disable
BROKER_URL=amqp://guest:guest@broker:5672/
COMMAND_SERVICE_PORT=8080
QUERY_SERVICE_PORT=8081
```

---

## Key Features Implemented

CQRS Pattern
Event-Driven Architecture
Outbox Pattern
RabbitMQ Integration
Docker Containerization
Materialized Views
Microservices Architecture
Fault-Tolerant System
Scalable Design

---

## Learning Outcomes

Understanding CQRS pattern
Understanding Event-Driven Systems
Using RabbitMQ message broker
Building scalable backend systems
Using Docker for microservices

---

## Conclusion

This project demonstrates a production-ready backend architecture using CQRS and event-driven design. It supports scalable analytics and ensures reliable communication between services.

This system is suitable for high-performance real-world applications.
