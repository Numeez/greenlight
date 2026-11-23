package data

import (
	"database/sql"
	"time"

	"github.com/Numeez/greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(validator *validator.Validator, movie *Movie) {
	validator.Check(movie.Title != "", "title", "must be provided")
	validator.Check(len(movie.Title) <= 500, "title", "must not be more than 500 characters")
	validator.Check(movie.Year != 0, "year", "must be provided")
	validator.Check(movie.Year > 1888, "year", "must be greater than 1888")
	validator.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	validator.Check(movie.Runtime != 0, "runtime", "must be provided")
	validator.Check(movie.Runtime > 0, "runtime", "must be positive integer")
	validator.Check(movie.Genres != nil, "genre", "must be provided")
	validator.Check(len(movie.Genres) >= 1, "genre", "must contain atleat 1 entry")
	validator.Check(len(movie.Genres) <= 5, "genre", "must not contain more than 5 entry")
	validator.Check(validator.Unique(movie.Genres), "genre", "must not contain duplicate values")

}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title,year,runtime,genres)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at,version
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

}

func (m MovieModel) Update(movie *Movie) (*Movie, error) {
	return nil, nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}
func (m MovieModel) Delete(id int64) error {
	return nil
}
