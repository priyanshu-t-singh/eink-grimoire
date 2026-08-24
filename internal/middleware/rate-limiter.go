package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimiter(next http.Handler) http.Handler {
	// 2 requests per second (rate.Limit(2)) with a max burst capacity of 5
	limiter := rate.NewLimiter(2, 5)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
