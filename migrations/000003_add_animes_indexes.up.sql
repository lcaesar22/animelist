CREATE INDEX IF NOT EXISTS animes_title_idx ON animes USING GIN (to_tsvector('simple', title));

CREATE INDEX IF NOT EXISTS animes_genre_idx ON animes USING GIN (to_tsvector(genre);