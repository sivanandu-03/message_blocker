-- WRITE MODEL TABLES

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    price NUMERIC NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    total NUMERIC NOT NULL,
    status VARCHAR(50) DEFAULT 'CREATED',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    price NUMERIC NOT NULL
);

-- OUTBOX TABLE (CRITICAL REQUIREMENT)

CREATE TABLE IF NOT EXISTS outbox (
    id SERIAL PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    published_at TIMESTAMP NULL
);

-- READ MODEL TABLES

CREATE TABLE IF NOT EXISTS product_sales_view (
    product_id INT PRIMARY KEY,
    total_quantity_sold INT DEFAULT 0,
    total_revenue NUMERIC DEFAULT 0,
    order_count INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS category_metrics_view (
    category_name VARCHAR(255) PRIMARY KEY,
    total_revenue NUMERIC DEFAULT 0,
    total_orders INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS customer_ltv_view (
    customer_id INT PRIMARY KEY,
    total_spent NUMERIC DEFAULT 0,
    order_count INT DEFAULT 0,
    last_order_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS hourly_sales_view (
    hour_timestamp TIMESTAMP PRIMARY KEY,
    total_orders INT DEFAULT 0,
    total_revenue NUMERIC DEFAULT 0
);

-- IDEMPOTENCY TABLE

CREATE TABLE IF NOT EXISTS processed_events (
    event_id INT PRIMARY KEY
);
