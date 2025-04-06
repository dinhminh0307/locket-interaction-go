package config

// CORSConfig defines CORS settings
type CORSConfig struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    AllowCredentials bool
    MaxAge           int // in seconds
}

// DefaultCORSConfig returns the default CORS configuration
func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"}, // Add your frontend origins here
        AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{
            "Authorization", "Content-Type", "Accept",
            "Origin", "User-Agent", "DNT", "Cache-Control",
            "X-Requested-With", "X-Client-Version", "X-Firebase-GMPID",
        },
        ExposedHeaders:   []string{},
        AllowCredentials: true,
        MaxAge:           86400, // 24 hours
    }
}