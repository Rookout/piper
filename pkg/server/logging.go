package server

import (
	"github.com/rs/zerolog"
	"net/http"
)

type LoggingMiddleware struct {
	logger *zerolog.Logger
}

func NewLoggingMiddleware(logger *zerolog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug(
			"request started",
			zap.String("proto", r.Proto),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}
