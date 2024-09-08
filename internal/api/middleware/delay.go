package middleware

import (
	"net/http"
	"time"
)

func Delay(next http.Handler, d time.Duration) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(d)
			next.ServeHTTP(w, r)
		},
	)
}
