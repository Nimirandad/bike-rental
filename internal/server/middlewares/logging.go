package middlewares

import (
	"net/http"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		log := logger.Get()

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Request started")

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		logEvent := log.Info()
		if wrapped.statusCode >= 400 {
			logEvent = log.Warn()
		}
		if wrapped.statusCode >= 500 {
			logEvent = log.Error()
		}

		logEvent.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.statusCode).
			Int("response_size", wrapped.written).
			Dur("duration_ms", duration).
			Msg("Request completed")
	})
}
