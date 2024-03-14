CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    login TEXT UNIQUE NOT NULL,
    role TEXT CONSTRAINT role_constratin CHECK (role = 'user' OR role = 'admin') DEFAULT 'user',
    password TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS login_idx ON users(login);

CREATE TABLE IF NOT EXISTS users_sessions (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id BIGSERIAL REFERENCES users(id),
    refresh_token TEXT NOT NULL,
    expired_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);