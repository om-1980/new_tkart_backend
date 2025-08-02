package middleware

import (
	"context"
	"net/http"
	"strings"

	"shared/jwt"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}