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
        ServerPort:    port,
        ServerHost:    getEnv("SERVER_HOST", ""),
		CORS:          DefaultCORSConfig(),
    }
    
    return cfg, nil
}

// Helper function to get environment variable or default.
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}