package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"locket-interaction-go/global"
	"locket-interaction-go/internal/models"
	"locket-interaction-go/internal/models/request"
	"locket-interaction-go/internal/models/response"
	"locket-interaction-go/pkg/firebase"
	"net/http"

	"time"
)

// Service handles authentication with the Locket API
type Service struct {
    firebaseService *firebase.Service
}

// NewService creates a new authentication service
func NewService(baseURL string, timeout time.Duration) (*Service, error) {
    apiKey := global.GetFirebaseAPIKey()
    if apiKey == "" {
        return nil, fmt.Errorf("Firebase API key not provided")
    }
    
    httpClient := &http.Client{
        Timeout: timeout,
    }
    
    firebaseService := firebase.NewService(httpClient, apiKey)
    
    return &Service{
        firebaseService: firebaseService,
    }, nil
}

// Login attempts to authenticate a user with the Locket API
func (s *Service) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
    // Call Firebase service to handle authentication
    firebaseResp, err := s.firebaseService.SignInWithEmailPassword(ctx, email, password)
    if err != nil {
        return nil, err
    }

    // Convert Firebase response to our application model
    loginResp, err := firebaseResp.ToLoginResponse()
    if err != nil {
        return nil, fmt.Errorf("failed to process Firebase response: %w", err)
    }

    return loginResp, nil
}

// LoginWithPhone attempts to authenticate a user with phone number and password
func (s *Service) LoginWithPhone(ctx context.Context, phone, password string) (*response.PhoneLoginResponse, error) {
    // Create the request body
    reqBody := request.PhoneLoginRequest{
        Data: request.PhoneLoginData{
            Phone:      phone,
            Password:   password,
            IOSVersion: global.IOSVersion,
        },
    }
    
    // Convert to JSON
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    // Create the request
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, global.LoginWithPhone, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    
    // Create HTTP client if needed
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Send the request
    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    // Check status code
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
    }
    
    // Parse the response
    var phoneLoginResp response.PhoneLoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&phoneLoginResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    return &phoneLoginResp, nil
}