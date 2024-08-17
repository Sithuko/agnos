CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    password VARCHAR(40) NOT NULL,
    steps INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
