package auth

import (
    "context"
    "fmt"
    "locket-interaction-go/global"
    "locket-interaction-go/internal/models"
    "locket-interaction-go/pkg"
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