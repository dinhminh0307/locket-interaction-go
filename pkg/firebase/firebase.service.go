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
	"time"

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

// UploadFileToStorage uploads a file to Firebase Storage
func (s *Service) UploadFileToStorage(ctx context.Context, userId, idToken string, fileData []byte, contentType string) (string, error) {
    // Generate a unique filename with timestamp
    fileName := fmt.Sprintf("%d_vtd182.webp", time.Now().UnixNano()/1000000) // Equivalent to Date.now()
    storagePath := fmt.Sprintf("users/%s/moments/thumbnails/%s", userId, fileName)
    encodedPath := url.PathEscape(storagePath)
    
    // Step 1: Initialize the upload
    initURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/locket-img/o/%s?uploadType=resumable&name=%s", 
                           encodedPath, encodedPath)
    
    // Create metadata for the upload
    metadata := map[string]interface{}{
        "name": storagePath,
        "contentType": contentType,
        "bucket": "",
        "metadata": map[string]string{
            "creator": userId,
            "visibility": "private",
        },
    }
    
    metadataJSON, err := json.Marshal(metadata)
    if err != nil {
        return "", fmt.Errorf("failed to marshal metadata: %w", err)
    }
    
    // Create request to initialize upload
    initReq, err := http.NewRequestWithContext(ctx, http.MethodPost, initURL, bytes.NewBuffer(metadataJSON))
    if err != nil {
        return "", fmt.Errorf("failed to create upload init request: %w", err)
    }
    
    // Set headers for initialization
    initReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
    initReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))
    initReq.Header.Set("X-Goog-Upload-Protocol", "resumable")
    initReq.Header.Set("Accept", "*/*")
    initReq.Header.Set("X-Goog-Upload-Command", "start")
    initReq.Header.Set("X-Goog-Upload-Content-Length", fmt.Sprintf("%d", len(fileData)))
    initReq.Header.Set("Accept-Language", "vi-VN,vi;q=0.9")
    initReq.Header.Set("X-Firebase-Storage-Version", "ios/10.13.0")
    initReq.Header.Set("User-Agent", "com.locket.Locket/1.43.1 iPhone/17.3 hw/iPhone15_3 (GTMSUF/1)")
    initReq.Header.Set("X-Goog-Upload-Content-Type", contentType)
    initReq.Header.Set("X-Firebase-Gmpid", "1:641029076083:ios:cc8eb46290d69b234fa609")
    
    // Send initialization request
    initResp, err := s.httpClient.Do(initReq)
    if err != nil {
        return "", fmt.Errorf("failed to send upload initialization request: %w", err)
    }
    defer initResp.Body.Close()
    
    if initResp.StatusCode != http.StatusOK {
        respBody, _ := io.ReadAll(initResp.Body)
        return "", fmt.Errorf("upload initialization failed with status %d: %s", 
                              initResp.StatusCode, string(respBody))
    }
    
    // Get the upload URL from the response headers
    uploadURL := initResp.Header.Get("X-Goog-Upload-URL")
    if uploadURL == "" {
        return "", fmt.Errorf("failed to get upload URL from response headers")
    }
    
    // Step 2: Upload the file data
    uploadReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, bytes.NewReader(fileData))
    if err != nil {
        return "", fmt.Errorf("failed to create upload request: %w", err)
    }
    
    // Add required headers for the upload request from global constants
    for key, value := range global.UPLOADER_HEADERS {
        uploadReq.Header.Set(key, value)
    }
    
    uploadResp, err := s.httpClient.Do(uploadReq)
    if err != nil {
        return "", fmt.Errorf("failed to upload file: %w", err)
    }
    defer uploadResp.Body.Close()
    
    if uploadResp.StatusCode != http.StatusOK {
        respBody, _ := io.ReadAll(uploadResp.Body)
        return "", fmt.Errorf("file upload failed with status %d: %s", 
                              uploadResp.StatusCode, string(respBody))
    }
    
    // Step 3: Get the download URL
    getURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/locket-img/o/%s", encodedPath)
    getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
    if err != nil {
        return "", fmt.Errorf("failed to create get token request: %w", err)
    }
    
    getReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
    getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))
    
    getResp, err := s.httpClient.Do(getReq)
    if err != nil {
        return "", fmt.Errorf("failed to get download token: %w", err)
    }
    defer getResp.Body.Close()
    
    if getResp.StatusCode != http.StatusOK {
        respBody, _ := io.ReadAll(getResp.Body)
        return "", fmt.Errorf("failed to get download token with status %d: %s", 
                              getResp.StatusCode, string(respBody))
    }
    
    // Parse the response to get the download token
    var respData map[string]interface{}
    respBody, err := io.ReadAll(getResp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %w", err)
    }
    
    if err := json.Unmarshal(respBody, &respData); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }
    
    downloadToken, ok := respData["downloadTokens"].(string)
    if !ok {
        return "", fmt.Errorf("failed to get download token from response")
    }
    
    // Return the complete download URL
    return fmt.Sprintf("%s?alt=media&token=%s", getURL, downloadToken), nil
}