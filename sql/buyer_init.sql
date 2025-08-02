
CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    mobile TEXT NOT NULL,
    address TEXT,
    district TEXT,
    state TEXT,
    country TEXT,
    pincode TEXT,
    profile_photo TEXT
);
