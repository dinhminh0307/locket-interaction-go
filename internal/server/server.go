package server

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "locket-interaction-go/config"
    "locket-interaction-go/internal/controllers"
    "locket-interaction-go/internal/models"
)

// Server represents the HTTP server for the application
type Server struct {
    router       *http.ServeMux
    httpServer   *http.Server
    locketClient *controllers.Client
    config       *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
    // Create Locket client
    locketClient, err := controllers.NewClient(cfg.LocketBaseURL, 30*time.Second)
    if err != nil {
        return nil, err
    }

    // Create server
    server := &Server{
        router:       http.NewServeMux(),
        locketClient: locketClient,
        config:       cfg,
    }

    // Register routes
    server.registerRoutes()

    // Configure HTTP server
    server.httpServer = &http.Server{
        Addr:         ":8080", // You could make this configurable
        Handler:      server.router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    return server, nil
}

// Start begins listening for requests
func (s *Server) Start() error {
    log.Printf("Starting server on %s", s.httpServer.Addr)
    return s.httpServer.ListenAndServe()
}

// registerRoutes sets up all the routes for the server
func (s *Server) registerRoutes() {
    // Register auth routes
    s.router.HandleFunc("/api/auth/login", s.handleLogin())
    
    // You can add more routes as needed
    // s.router.HandleFunc("/api/some-endpoint", s.handleSomeEndpoint())
}

// handleLogin returns an HTTP handler for the login endpoint
func (s *Server) handleLogin() http.HandlerFunc {
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

        // Call Locket API through our client
        loginResp, err := s.locketClient.Login(ctx, loginReq.Email, loginReq.Password)
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