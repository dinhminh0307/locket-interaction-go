package middleware

import (
    "net/http"
    "strings"

    "locket-interaction-go/config"
)

// CORSMiddleware applies CORS headers to responses
func CORSMiddleware(corsConfig config.CORSConfig) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if the request origin is allowed
            allowed := false
            for _, allowedOrigin := range corsConfig.AllowedOrigins {
                if allowedOrigin == "*" || origin == allowedOrigin {
                    allowed = true
                    break
                }
            }

            if allowed {
                // Set CORS headers
                w.Header().Set("Access-Control-Allow-Origin", origin)
                
                if corsConfig.AllowCredentials {
                    w.Header().Set("Access-Control-Allow-Credentials", "true")
                }

                // Handle preflight requests
                if r.Method == http.MethodOptions {
                    w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
                    w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
                    
                    if corsConfig.MaxAge > 0 {
                        w.Header().Set("Access-Control-Max-Age", string(corsConfig.MaxAge))
                    }
                    
                    w.WriteHeader(http.StatusOK)
                    return
                }
            }

            next.ServeHTTP(w, r)
        })
    }
}