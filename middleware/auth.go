package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AshimKoirala/load-balancer-admin/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("token")

		token, err := extractBearerToken(r)

		if err != nil {
			// if err == http.ErrNoCookie {
			// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// 	return
			// }
			http.Error(w, "authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(token)

		if err != nil {
			log.Print(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		// Set the username in the context
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) (string, error) {
	// Retrieve the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Check if the header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("authorization header is not a bearer token")
	}

	// Extract the token by trimming the prefix
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	return token, nil
}
