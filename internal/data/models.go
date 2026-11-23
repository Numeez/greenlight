package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Movies interface {
	Insert(movie *Movie) error
	Get(id int64) (*Movie, error)
	Update(movie *Movie) (*Movie, error)
	Delete(id int64) error
}

type Model struct {
	Movies Movies
}

func NewModels(db *sql.DB) Model {
	return Model{
		Movies: MovieModel{
			DB: db,
		},
	}
}
