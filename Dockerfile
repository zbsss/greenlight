FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server ./cmd/api

FROM alpine:latest

COPY --from=build /server /app/server

CMD ["/app/server"]

EXPOSE 4000
