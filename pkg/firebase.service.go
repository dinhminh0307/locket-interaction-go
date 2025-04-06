package firebase

import (
    "bytes"
    "compress/flate"
    "compress/gzip"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"

    "locket-interaction-go/global"
    "locket-interaction-go/internal/models"
)

// Service provides methods for interacting with Firebase Authentication
type Service struct {
    httpClient *http.Client
    apiKey     string
}

// NewService creates a new Firebase service instance
func NewService(httpClient *http.Client, apiKey string) *Service {
    return &Service{
        httpClient: httpClient,
        apiKey:     apiKey,
    }
}

// SignInWithEmailPassword authenticates a user with Firebase using email/password
func (s *Service) SignInWithEmailPassword(ctx context.Context, email, password string) (*models.FirebaseAuthResponse, error) {
    // 1. Prepare Request Body
    reqBody := models.LoginRequest{
        Email:            email,
        Password:         password,
        ReturnSecureToken: true,
        ClientType:       "CLIENT_TYPE_IOS",
    }
    
    return s.sendAuthRequest(ctx, global.LoginURL, reqBody)
}

// sendAuthRequest sends a request to Firebase Auth API
func (s *Service) sendAuthRequest(ctx context.Context, endpoint string, reqBody interface{}) (*models.FirebaseAuthResponse, error) {
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Construct URL with API key
    endpointURL, err := url.Parse(endpoint)
    if err != nil {
        return nil, fmt.Errorf("invalid endpoint URL: %w", err)
    }
    
    // Add API key to query string
    q := endpointURL.Query()
    q.Add(global.APIKeyQueryParam, s.apiKey)
    endpointURL.RawQuery = q.Encode()

    // Create Request
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL.String(), bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Add required headers
    for key, value := range global.LoginHeaders {
        req.Header.Set(key, value)
    }
    
    // Add Firebase-specific auth headers
    for key, value := range global.FirebaseAuthHeaders {
        req.Header.Set(key, value)
    }

    // Execute Request
    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Handle Response
    respBody, err := s.readResponseBody(resp)
    if err != nil {
        return nil, err
    }

    // Check for non-success status codes
    if resp.StatusCode != http.StatusOK {
        return nil, s.handleErrorResponse(resp, respBody)
    }

    // Parse successful response
    var firebaseResp models.FirebaseAuthResponse
    if err := json.Unmarshal(respBody, &firebaseResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(respBody))
    }

    return &firebaseResp, nil
}

// readResponseBody reads and decompresses the response body if needed
func (s *Service) readResponseBody(resp *http.Response) ([]byte, error) {
    var reader io.ReadCloser
    var err error

    // Handle compressed responses
    switch resp.Header.Get("Content-Encoding") {
    case "gzip":
        reader, err = gzip.NewReader(resp.Body)
        if err != nil {
            return nil, fmt.Errorf("failed to create gzip reader: %w", err)
        }
        defer reader.Close()
    case "deflate":
        reader = flate.NewReader(resp.Body)
        defer reader.Close()
    default:
        reader = resp.Body
    }

    return io.ReadAll(reader)
}

// handleErrorResponse processes error responses from Firebase
func (s *Service) handleErrorResponse(resp *http.Response, body []byte) error {
    // Log the raw bytes for debugging
    log.Printf("Raw error response (%d bytes): %v", len(body), body)
    
    // Try to decode as JSON first
    var prettyJSON bytes.Buffer
    if json.Indent(&prettyJSON, body, "", "  ") == nil {
        // Check if it's a Firebase error
        var firebaseError models.FirebaseErrorResponse
        if err := json.Unmarshal(body, &firebaseError); err == nil {
            return firebaseError.ToErrorResponse()
        }
        
        // It's valid JSON but not in Firebase error format
        return fmt.Errorf("request failed: status code %d, response: %s", 
            resp.StatusCode, prettyJSON.String())
    }
    
    // If it's not JSON, return as text with hex dump for binary data
    return fmt.Errorf("request failed: status code %d, response type: %s, body: %s, hex: %x",
        resp.StatusCode, resp.Header.Get("Content-Type"), string(body), body)
}