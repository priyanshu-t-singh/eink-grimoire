package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, x := range slices.Backward(xs) {
			next = x(next)
		}

		return next
	}
}
