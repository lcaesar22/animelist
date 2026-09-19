package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/konaha-gakure/otaku/internal/validator"
	"github.com/lib/pq"
)

type Anime struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero,string"`
	Genre     []string  `json:"genre"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	Version   int       `json:"version"`
}

// anime model
type AnimeModel struct {
	DB *sql.DB
}

func (m AnimeModel) GetAll(title string, genre []string, filters Filters) ([]*Anime, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, title, year, genre, created_at, version FROM animes 
                                                   WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
                                                              AND (genre @> $2 OR $2 = '{}') ORDER BY %s %s, id ASC 
                                                              LIMIT $3  OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genre), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalPage := 0
	animes := []*Anime{}

	for rows.Next() {
		var anime Anime

		err := rows.Scan(
			&totalPage,
			&anime.ID,
			&anime.Title,
			&anime.Year,
			&anime.Runtime,
			pq.Array(&anime.Genre),
			&anime.CreatedAt,
			&anime.Version)

		if err != nil {
			return nil, Metadata{}, err
		}

		animes = append(animes, &anime)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalPage, filters.Page, filters.PageSize)

	return animes, metadata, nil
}

// placeholder for inserting a new record
func (m AnimeModel) Insert(anime *Anime) error {
	query := `INSERT INTO anime (title, year, runtime, genre)
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at, updated_at. version`

	args := []interface{}{anime.Title, anime.Year, anime.Runtime, pq.Array(anime.Genre)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.ID)
}

// ” for fetching a specific record
func (m AnimeModel) Get(id int) (*Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// sql query
	query := `SELECT id, title, year, runtime, genre, created_at, version FROM animes WHERE id = $1`

	var anime Anime

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&anime.ID,
		&anime.Title,
		&anime.Year,
		&anime.Runtime,
		pq.Array(&anime.Genre),
		&anime.CreatedAt,
		pq.Array(&anime.Version))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &anime, nil
}

// ” for updates
func (m AnimeModel) Update(anime *Anime) error {
	query := `UPDATE animes SET title = %1, year = %2, runtime = %3, genre = %4, version = uuid_generate_v4() WHERE id = $5 AND version = $6 RETURNING version`

	args := []interface{}{
		anime.Title,
		anime.Year,
		anime.Runtime,
		pq.Array(anime.Genre),
		anime.Version,
		anime.ID,
		anime.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// triggers optimistic licking
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	// execute the query
	return nil
}

// ” for deletion
func (m AnimeModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// query. delete record
	query := `DELETE FROM animes WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// triggers if no rows are affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// validation code
func ValidateAnime(v *validator.Validator, anime *Anime) {

	v.Check(anime.Title != "", "title", "required")
	v.Check(len(anime.Title) <= 500, "title", "must not be more than 500 bytes")

	v.Check(anime.Year != 0, "year", "year must be provided")
	v.Check(anime.Year >= 1888, "year", "must be greater or equal to 1888")
	v.Check(anime.Year <= int(time.Now().Year()), "year", "must not be in the future")

	v.Check(anime.Runtime != 0, "runtime", "must be provided")
	v.Check(anime.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(anime.Genre != nil, "genre", "must be provided")
	v.Check(len(anime.Genre) >= 1, "genre", "must contain at least 1 genre")
	v.Check(len(anime.Genre) <= 5, "genre", "must contain at most 5 genres")

	// check if the slices are unique
	v.Check(validator.Unique(anime.Genre), "genre", "must not contain duplicate values")
}
