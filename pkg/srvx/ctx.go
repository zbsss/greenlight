package srvx

type ctxKey string

const (
	traceIDKey       ctxKey = "traceID"
	requestLoggerKey ctxKey = "requestLogger"
)
