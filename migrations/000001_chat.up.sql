CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    user VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    timestamp BIGINT NOT NULL,
    ip_address VARCHAR(255) NOT NULL
);
