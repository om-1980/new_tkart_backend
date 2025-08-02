
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    delivery_id VARCHAR(20) NOT NULL UNIQUE,
    order_id INT NOT NULL,
    delivery_person VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    is_cod BOOLEAN DEFAULT FALSE,
    cod_amount INT DEFAULT 0,
    is_return BOOLEAN DEFAULT FALSE
);
