## Prerequisites

```
brew install sqlc golang-migrate

go install github.com/air-verse/air@latest github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Database migrations

### Create a new migration

```sh
migrate create -seq -ext=.sql -dir=./movies/backend/storage/migrations create_movies_table
```

### Run migrations

```sh
migrate -path ./movies/backend/storage/migrations -database "postgres://greenlight:password@localhost:5432/greenlight?sslmode=disable" up
```

### Generate Go models from SQL schema

```sh
sqlc generate -f movies/backend/storage/sqlc.yaml
```

## Generate API

### Generate server

```sh
oapi-codegen -generate types,std-http-server -package api -o movies/backend/api/openapi.gen.go movies/api/movies.yaml
```

### Generate client

```sh
cd movies/frontend
npx openapi-typescript ../api/movies.yaml -o src/lib/api/v1.d.ts
```
