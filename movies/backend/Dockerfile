FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY movies/backend/ ./movies/backend/

RUN go build -o /movies-server ./movies/backend

FROM alpine:3.21.3

# Create a non-root user to run the application
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

# Create directory for the application and set proper ownership
RUN mkdir -p /app && chown -R nonroot:nonroot /app

WORKDIR /app

COPY --from=build /movies-server /app/movies-server

# Switch to non-root user
USER nonroot

CMD ["/app/movies-server"]

EXPOSE 4000
