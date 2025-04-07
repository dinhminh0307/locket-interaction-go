package config

import (
    "os"
    "strconv"
)

// Config holds the application configuration.
type Config struct {
    LocketBaseURL string
    LocketAPIKey  string
    ServerPort    int
    ServerHost    string
    FirebaseAPIKey string
    CreatePostURL  string
	CORS          CORSConfig
}

// Load loads configuration from environment variables or defaults.
func Load() (*Config, error) {
    port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        port = 8080
    }

    // Base URL for Locket API
    baseURL := getEnv("LOCKET_BASE_URL", "https://www.googleapis.com/identitytoolkit/v3/relyingparty")
    
    cfg := &Config{
        LocketBaseURL: baseURL,
        LocketAPIKey:  getEnv("LOCKET_API_KEY", ""),
        CreatePostURL:  "https://api.locketcamera.com/postMomentV2",
        FirebaseAPIKey: getEnv("FIREBASE_API_KEY", ""),
        ServerPort:    port,
        ServerHost:    getEnv("SERVER_HOST", ""),
		CORS:          DefaultCORSConfig(),
    }
    
    return cfg, nil
}

// return &Config{
//     ServerPort:     8080, // Replace with actual loaded value
//     LocketBaseURL:  "https://api.locket.com", // Replace with actual loaded value
//     FirebaseAPIKey: "your-firebase-api-key", // Replace with actual loaded value
//     CreatePostURL:  "https://api.locketcamera.com/postMomentV2", // Replace with actual loaded value
//     CORS: CORSConfig{
//         AllowedOrigins: []string{"*"}, // Replace with actual loaded value
//         AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Replace with actual loaded value
//         AllowedHeaders: []string{"Content-Type", "Authorization"}, // Replace with actual loaded value
//     },
// }, nil

// Helper function to get environment variable or default.
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}