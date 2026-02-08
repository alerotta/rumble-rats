-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users;

CREATE TABLE users (
  username       TEXT PRIMARY KEY,
  email          CITEXT NOT NULL UNIQUE,
  password_hash  TEXT   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;