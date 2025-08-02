CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100),
    mobile VARCHAR(15),
    name VARCHAR(100),
    address TEXT,
    district VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    pincode VARCHAR(10),
    date TIMESTAMP,
    status VARCHAR(20) DEFAULT 'Placed'
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    name VARCHAR(100),
    price NUMERIC(10, 2),
    qty INTEGER,
    image TEXT,
    seller_id VARCHAR(20)
);
