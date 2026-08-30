-- Initial schema for the tennis-court booking bot.
-- This template ships no business logic; adapt these tables to real requirements.

CREATE TABLE IF NOT EXISTS courts (
    id     BIGSERIAL PRIMARY KEY,
    name   TEXT    NOT NULL,
    indoor BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT      NOT NULL UNIQUE,
    username    TEXT        NOT NULL DEFAULT '',
    full_name   TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bookings (
    id         BIGSERIAL PRIMARY KEY,
    court_id   BIGINT      NOT NULL REFERENCES courts (id),
    user_id    BIGINT      NOT NULL REFERENCES users (id),
    starts_at  TIMESTAMPTZ NOT NULL,
    ends_at    TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS bookings_user_id_idx ON bookings (user_id);
CREATE INDEX IF NOT EXISTS bookings_court_starts_at_idx ON bookings (court_id, starts_at);
