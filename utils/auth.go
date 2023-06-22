package utils

import (
	"net/http"
	"strings"
)

// Middleware to authenticate requests using JWT token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from request header
		token := extractToken(r)

		// Validate and verify the token
		if !validateToken(token) {
			SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// If token is valid, proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// Function to extract JWT token from request header
func extractToken(r *http.Request) string {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Remove "Bearer " prefix from token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token
}

// Function to validate and verify the JWT token
func validateToken(token string) bool {
	// Implement token validation logic
	// Return true if token is valid, false otherwise
	return true
}
