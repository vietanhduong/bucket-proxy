package server

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Errorf("Recover panic: %v Trace: \n%s", err, debug.Stack())
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			if wrapped.status == 0 {
				wrapped.status = http.StatusOK
			}
			// "POST /api/ingester.v1.PushService/Push HTTP/1.1" 400 34 229.301µs "127.0.0.1:41046" "curl/7.68.0" "anhdv" "localhost:8080"
			log.WithFields(logrus.Fields{
				"http.method":         r.Method,
				"http.url":            r.URL.EscapedPath(),
				"http.proto":          r.Proto,
				"http.status":         wrapped.status,
				"http.content-length": r.ContentLength,
				"http.duration":       time.Since(start),
				"http.remote-addr":    extractRemoteAddr(r),
				"http.user-agent":     r.Header.Get("user-agent"),
				"http.request-id":     r.Header.Get("x-request-id"),
				"http.host":           r.Host,
			}).Info()
		}

		return http.HandlerFunc(fn)
	}
}

func extractRemoteAddr(r *http.Request) string {
	if val := r.Header.Get("x-forwarded-for"); val != "" {
		return val
	}
	return r.RemoteAddr
}
