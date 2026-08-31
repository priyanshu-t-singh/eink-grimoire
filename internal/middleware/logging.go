package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

const logContextKey contextKey = "middleware.logging.contextHolder"

type ContextHolder struct {
	Ctx context.Context
}

func UpdateLoggingContext(r *http.Request, ctx context.Context) {
	if holder, ok := r.Context().Value(logContextKey).(*ContextHolder); ok {
		holder.Ctx = ctx
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		holder := &ContextHolder{Ctx: r.Context()}
		r = r.WithContext(context.WithValue(r.Context(), logContextKey, holder))

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		slog.InfoContext(
			holder.Ctx,
			"request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}
