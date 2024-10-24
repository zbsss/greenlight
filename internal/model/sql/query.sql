-- name: ListMovies :many
SELECT * FROM movies;

-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetMovie :one
SELECT * FROM movies
WHERE id = $1;
