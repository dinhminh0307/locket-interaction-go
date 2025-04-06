package controllers

import (
    "bytes"
    "compress/flate"
    "compress/gzip"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "locket-interaction-go/global"
    "locket-interaction-go/internal/models"
    "net/http"
    "net/url"
    "time"
)

// Client manages communication with the Locket backend API.
type Client struct {
    httpClient *http.Client
    baseURL    *url.URL
    apiKey     string
}

// NewClient creates a new Locket API client.
// NewClient creates a new Locket API client.
func NewClient(baseURL string, timeout time.Duration) (*Client, error) {
    parsedBaseURL, err := url.Parse(baseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid base URL: %w", err)
    }

    apiKey := global.GetFirebaseAPIKey()
    if apiKey == "" {
        return nil, fmt.Errorf("Firebase API key not provided")
    }

    return &Client{
        httpClient: &http.Client{
            Timeout: timeout,
        },
        baseURL: parsedBaseURL,
        apiKey:  apiKey,
    }, nil
}

// Update the Login method in the controllers package
func (c *Client) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
     // 1. Prepare Request Body
	 reqBody := models.LoginRequest{
        Email:            email,
        Password:         password,
        ReturnSecureToken: true,
        ClientType:       "CLIENT_TYPE_IOS",
    }
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal login request: %w", err)
    }

    // 2. Construct URL
    loginURL, err := url.Parse(global.LoginURL)
    if err != nil {
        return nil, fmt.Errorf("invalid login URL: %w", err)
    }
    
    // Add API key to query string - FIXED: Proper query parameter construction
    q := loginURL.Query()
    q.Add(global.APIKeyQueryParam, c.apiKey) // Add(key, value)
    loginURL.RawQuery = q.Encode()

    // 3. Create Request
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL.String(), bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create login request: %w", err)
    }
    
    // Add required headers
    for key, value := range global.LoginHeaders {
        req.Header.Set(key, value)
    }
    
    // Add Firebase-specific auth headers
    for key, value := range global.FirebaseAuthHeaders {
        req.Header.Set(key, value)
    }

    // 4. Execute Request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute login request: %w", err)
    }
    defer resp.Body.Close()

    // 5. Handle Response
    var reader io.ReadCloser

	// Check if the response is compressed
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

	// Read the response body
	respBodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response body: %w", err)
	}

	// Check for non-success status codes
	if resp.StatusCode != http.StatusOK {
		// Log the raw bytes for debugging
		log.Printf("Raw response (%d bytes): %v", len(respBodyBytes), respBodyBytes)
		
		// Try to decode as JSON first
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, respBodyBytes, "", "  ") == nil {
			// It's valid JSON, return it formatted nicely
			return nil, fmt.Errorf("login failed: status code %d, response: %s", 
				resp.StatusCode, prettyJSON.String())
		}
		
		// Attempt to parse as Firebase error format
		var firebaseError models.FirebaseErrorResponse
		if err := json.Unmarshal(respBodyBytes, &firebaseError); err == nil {
			return nil, firebaseError.ToErrorResponse()
		}
		
		// If it's not JSON, return as text with hex dump for binary data
		return nil, fmt.Errorf("login failed: status code %d, response type: %s, body: %s, hex: %x",
			resp.StatusCode, resp.Header.Get("Content-Type"), string(respBodyBytes), respBodyBytes)
	}

    // Parse successful response as Firebase auth response
    var firebaseResp models.FirebaseAuthResponse
    if err := json.Unmarshal(respBodyBytes, &firebaseResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal Firebase response: %w, body: %s", err, string(respBodyBytes))
    }

    // Convert to our standard format
    loginResp, err := firebaseResp.ToLoginResponse()
    if err != nil {
        return nil, fmt.Errorf("failed to process Firebase response: %w", err)
    }

    return loginResp, nil
}