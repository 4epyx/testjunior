package middleware

import (
	"net/http"
	"time"

	hu "github.com/4epyx/testtask/util/http"
	"github.com/rs/zerolog"
)

type LogMiddleware struct {
	l *zerolog.Logger
}

func NewLogMiddleware(logger *zerolog.Logger) *LogMiddleware {

	return &LogMiddleware{l: logger}
}

func (m *LogMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := hu.WrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		end := time.Now()

		if wrapped.Status() != 500 {
			m.l.Info().
				Str("method", r.Method).
				Time("timestamp", start).
				Str("path", r.URL.Path).
				Dur("latency", end.Sub(start)).
				Int("status", wrapped.Status()).
				Send()
		} else {
			m.l.Error().
				Str("method", r.Method).
				Time("timestamp", start).
				Str("path", r.URL.Path).
				Dur("latency", end.Sub(start)).
				Int("status", wrapped.Status()).
				Str("msg", wrapped.StringMsg()).
				Send()
		}
	})
}
