## Database migrations

### Install migrate

```sh
brew install golang-migrate
```

### Create a new migration

```sh
migrate create -seq -ext=.sql -dir=./internal/movies/model/sql/migrations create_movies_table
```

### Run migrations

```sh
migrate -path ./internal/movies/model/sql/migrations -database "postgres://greenlight:password@localhost:5432/greenlight?sslmode=disable" up
```

### Generate Go models from SQL schema

```sh
sqlc generate
```
