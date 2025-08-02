CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    seller_id VARCHAR(20) NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    subcategory TEXT NOT NULL,
    inner_subcategory TEXT NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    in_stock BOOLEAN NOT NULL,
    image1 TEXT,
    image2 TEXT,
    image3 TEXT,
    image4 TEXT
);
