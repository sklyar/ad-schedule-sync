package http

import "net/http"

const secretKeyHeader = "X-Secret-Key"

// checkSecretKey is a middleware that checks if the request has a valid secret key.
func checkSecretKeyMiddleware(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get(secretKeyHeader)
			if key != secretKey {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
