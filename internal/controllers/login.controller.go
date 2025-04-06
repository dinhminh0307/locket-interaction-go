package controllers

import (
    "context"
    "encoding/json"
    "log"
    "locket-interaction-go/internal/models"
    "locket-interaction-go/internal/services/auth"
    "net/http"
    "time"
)

// AuthController handles HTTP requests related to authentication
type AuthController struct {
    authService *auth.Service
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService *auth.Service) *AuthController {
    return &AuthController{
        authService: authService,
    }
}

// HandleLogin handles login requests
func (c *AuthController) HandleLogin() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Only allow POST requests
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Parse request body
        var loginReq models.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate required fields
        if loginReq.Email == "" || loginReq.Password == "" {
            http.Error(w, "Email and password are required", http.StatusBadRequest)
            return
        }

        // Create context with timeout
        ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
        defer cancel()

        // Call auth service to handle login
        loginResp, err := c.authService.Login(ctx, loginReq.Email, loginReq.Password)
        if err != nil {
            // Check if it's our custom error type
            if apiErr, ok := err.(*models.ErrorResponse); ok {
                w.WriteHeader(apiErr.Code)
                json.NewEncoder(w).Encode(apiErr)
                return
            }
            
            // Generic error
            log.Printf("Login error: %v", err)
            http.Error(w, "Authentication failed", http.StatusInternalServerError)
            return
        }

        // Return successful response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(loginResp)
    }
}