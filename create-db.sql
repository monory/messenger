-- PostgreSQL

CREATE DATABASE messenger;
\connect messenger;
CREATE USER messenger_user WITH PASSWORD 'example_password';
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO messenger_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO messenger_user;


CREATE TABLE Users (
    id              BIGSERIAL PRIMARY KEY,
    login           VARCHAR(64) NOT NULL,
    password_hash   BYTEA NOT NULL,
    shown_name      VARCHAR(256),
    online          BOOLEAN
);

CREATE TABLE Chats (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(256) NOT NULL
);

CREATE TABLE Messages (
    id              BIGSERIAL PRIMARY KEY,
    time_stamp      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message         VARCHAR(10000),
    author_id       BIGINT NOT NULL REFERENCES Users(id),
    chat_id         BIGINT NOT NULL REFERENCES Chats(id),
    seen            BOOLEAN DEFAULT FALSE
);

CREATE TABLE UserContacts (
    id              BIGSERIAL PRIMARY KEY,
    owner_id        BIGINT NOT NULL REFERENCES Users(id),
    user_id         BIGINT NOT NULL REFERENCES Users(id),
    pseudonym       VARCHAR(256)
);

CREATE TABLE ChatContacts (
	id              BIGSERIAL PRIMARY KEY,
	owner_id        BIGINT NOT NULL REFERENCES Users(id),
	chat_id         BIGINT NOT NULL REFERENCES Chats(id),
	pseudonym       VARCHAR(256)
);

CREATE TABLE Tokens (
    id              SERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES Users(id),
    token           BYTEA
);
