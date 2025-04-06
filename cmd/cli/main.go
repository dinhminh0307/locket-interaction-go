package main

import (
    "log"

    "locket-interaction-go/config"
    "locket-interaction-go/global"
    "locket-interaction-go/internal/server"
)

func main() {
    // Initialize environment variables from .env
    global.Initialize()
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create and start server
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    // Start listening for requests
    log.Fatal(srv.Start())
}