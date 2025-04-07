package middleware

import (
	"context"
	"log"  // Added log package
	"net/http"
	"strings"
)

// RequireAuth is a middleware that checks for valid authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Auth failed: No Authorization header received")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Auth failed: Invalid authorization format")
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Log token reception (only showing first few chars for security)
		tokenPreview := idToken
		if len(idToken) > 10 {
			tokenPreview = idToken[:10] + "..."
		}
		log.Printf("Auth token received: %s", tokenPreview)
		
		// handle user id in header
		userHeader := r.Header.Get("X-User-ID")
		if userHeader == "" {
			log.Println("Auth failed: No userHeader header received")
			http.Error(w, "userHeader header required", http.StatusUnauthorized)
			return
		}

		userID := userHeader // Replace with actual user ID extraction
		
		ctx := context.WithValue(r.Context(), "userId", userID)
		ctx = context.WithValue(ctx, "idToken", idToken)

		// Call the next handler with the enriched context
		next(w, r.WithContext(ctx))
	}
}