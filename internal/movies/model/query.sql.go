// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package model

import (
	"context"
)

const createMovie = `-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4) RETURNING id, created_at, title, year, runtime, genres, version
`

type CreateMovieParams struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

func (q *Queries) CreateMovie(ctx context.Context, arg CreateMovieParams) (Movie, error) {
	row := q.db.QueryRow(ctx, createMovie,
		arg.Title,
		arg.Year,
		arg.Runtime,
		arg.Genres,
	)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Title,
		&i.Year,
		&i.Runtime,
		&i.Genres,
		&i.Version,
	)
	return i, err
}

const getMovie = `-- name: GetMovie :one
SELECT id, created_at, title, year, runtime, genres, version FROM movies
WHERE id = $1
`

func (q *Queries) GetMovie(ctx context.Context, id int64) (Movie, error) {
	row := q.db.QueryRow(ctx, getMovie, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Title,
		&i.Year,
		&i.Runtime,
		&i.Genres,
		&i.Version,
	)
	return i, err
}

const listMovies = `-- name: ListMovies :many
SELECT id, created_at, title, year, runtime, genres, version FROM movies
`

func (q *Queries) ListMovies(ctx context.Context) ([]Movie, error) {
	rows, err := q.db.Query(ctx, listMovies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Movie
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Title,
			&i.Year,
			&i.Runtime,
			&i.Genres,
			&i.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}