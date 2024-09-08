package middleware

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler

func NewHandlerWithMiddleware(h http.Handler, m ...MiddlewareFunc) http.Handler {
	current := h
	for _, v := range m {
		current = v(current)
	}

	return current
}
