CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email ci text UNIQUE NOT NULL,
    password_hash bytes NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);