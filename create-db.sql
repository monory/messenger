-- PostgreSQL

CREATE DATABASE messenger;
\connect messenger;

CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(64) UNIQUE NOT NULL,
    password_hash   BYTEA NOT NULL
);

CREATE TABLE messages (
    id              BIGSERIAL PRIMARY KEY,
    time_stamp      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message         VARCHAR(10000) NOT NULL,
    author_id       BIGINT NOT NULL REFERENCES users(id),
    receiver_id     BIGINT NOT NULL REFERENCES users(id),
    seen            BOOLEAN DEFAULT FALSE
);

CREATE TABLE contacts (
    id              BIGSERIAL PRIMARY KEY,
    owner_id        BIGINT NOT NULL REFERENCES users(id),
    user_id         BIGINT NOT NULL REFERENCES users(id)
);

CREATE TABLE tokens (
    id              SERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    selector        BYTEA NOT NULL UNIQUE,
    token           BYTEA NOT NULL,
    expires         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 week'
);

CREATE TABLE chat_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES Users(id),
    selector        BYTEA NOT NULL UNIQUE,
    token           BYTEA NOT NULL,
    expires         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour'
);

CREATE INDEX ON tokens USING HASH (selector);
CREATE INDEX ON chat_tokens USING HASH (selector);

CREATE USER web_backend WITH PASSWORD 'web_backend_password';
GRANT SELECT, INSERT ON TABLE users, tokens, chat_tokens TO web_backend;
GRANT ALL ON SEQUENCE users_id_seq, tokens_id_seq, chat_tokens_id_seq TO web_backend;

CREATE USER chat_backend WITH PASSWORD 'chat_backend_password';
GRANT ALL ON ALL TABLES IN SCHEMA public TO chat_backend;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO chat_backend;
