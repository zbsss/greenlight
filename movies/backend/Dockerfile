FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /movies-server ./movies/backend

FROM alpine:latest

COPY --from=build /movies-server /app/movies-server

CMD ["/app/movies-server"]

EXPOSE 4000
