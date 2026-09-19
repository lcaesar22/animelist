CREATE TABLE IF NOT EXISTS animes(
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone DEFAULT NOW()
);