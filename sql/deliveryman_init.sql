CREATE TABLE IF NOT EXISTS deliveryman (
    id SERIAL PRIMARY KEY,
    delivery_id VARCHAR(20) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    mobile TEXT NOT NULL,
    account_number TEXT,
    address TEXT,
    district TEXT,
    state TEXT,
    country TEXT,
    pincode TEXT,
    profile_photo TEXT,
    is_active BOOLEAN DEFAULT TRUE
)