CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS session_token_table(
    username VARCHAR(255) NOT NULL UNIQUE,
    session_token VARCHAR(255) NOT NULL UNIQUE
);