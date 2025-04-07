package server

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "locket-interaction-go/config"
    "locket-interaction-go/internal/controllers"
    "locket-interaction-go/internal/middlewares"
    "locket-interaction-go/internal/services/auth"
    "locket-interaction-go/internal/services/upload"
    "locket-interaction-go/pkg/firebase"
)

// Server represents the HTTP server for the application
type Server struct {
    router          *http.ServeMux
    httpServer      *http.Server
    authService     *auth.Service
    authController  *controllers.AuthController
    uploadService   *upload.Service
    uploadController *controllers.UploadController
    config          *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
    // Create HTTP client
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Create auth service
    authService, err := auth.NewService(cfg.LocketBaseURL, 30*time.Second)
    if err != nil {
        return nil, err
    }

    // Create Firebase service
    firebaseService := firebase.NewService(httpClient, cfg.FirebaseAPIKey)
    
    // Create upload service
    uploadService := upload.NewService(
        firebaseService,
        httpClient,
        cfg.CreatePostURL, // You need to add this to your config
    )

    // Create controllers
    authController := controllers.NewAuthController(authService)
    uploadController := controllers.NewUploadController(uploadService)

    // Create server
    server := &Server{
        router:          http.NewServeMux(),
        authService:     authService,
        authController:  authController,
        uploadService:   uploadService,
        uploadController: uploadController,
        config:          cfg,
    }

    // Register routes
    server.registerRoutes()

    // Apply CORS middleware
    corsHandler := middleware.CORSMiddleware(cfg.CORS)(server.router)

    // Configure HTTP server
    server.httpServer = &http.Server{
        Addr:         ":" + fmt.Sprintf("%d", cfg.ServerPort),
        Handler:      corsHandler,
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
    s.router.HandleFunc("/api/auth/login", s.authController.HandleLogin())
    
    // Register upload routes
    s.router.HandleFunc("/api/upload/image", s.uploadController.HandleUploadImage())
}