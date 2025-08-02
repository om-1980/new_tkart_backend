package middleware

import (
	"net/http"
	"strings"

	"shared/jwt"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		_, err := jwt.ParseJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
