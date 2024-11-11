-- name: ListMovies :many
SELECT * FROM movies;

-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime_min, genres)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetMovie :one
SELECT * FROM movies
WHERE id = $1;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2, year = $3, runtime_min = $4, genres = $5, version = version + 1
WHERE id = $1
RETURNING *;
