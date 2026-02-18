-- WRITE MODEL
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    category VARCHAR(255),
    price NUMERIC(10, 2),
    stock INTEGER
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER,
    total NUMERIC(10, 2),
    status VARCHAR(50) DEFAULT 'CREATED',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER,
    price NUMERIC(10, 2)
);

-- TRANSACTIONAL OUTBOX
CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP NULL
);

-- READ MODEL (Materialized Views)
CREATE TABLE product_sales_view (
    product_id INTEGER PRIMARY KEY,
    total_quantity_sold INTEGER DEFAULT 0,
    total_revenue NUMERIC(12, 2) DEFAULT 0,
    order_count INTEGER DEFAULT 0
);

CREATE TABLE category_metrics_view (
    category_name VARCHAR(255) PRIMARY KEY,
    total_revenue NUMERIC(12, 2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0
);

CREATE TABLE customer_ltv_view (
    customer_id INTEGER PRIMARY KEY,
    total_spent NUMERIC(12, 2) DEFAULT 0,
    order_count INTEGER DEFAULT 0,
    last_order_date TIMESTAMP
);

CREATE TABLE hourly_sales_view (
    hour_timestamp TIMESTAMP PRIMARY KEY,
    total_orders INTEGER DEFAULT 0,
    total_revenue NUMERIC(12, 2) DEFAULT 0
);

CREATE TABLE sync_status (
    id INTEGER PRIMARY KEY,
    last_event_time TIMESTAMP
);

INSERT INTO sync_status (id, last_event_time) VALUES (1, NOW());