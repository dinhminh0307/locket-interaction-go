package upload

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "locket-interaction-go/pkg/firebase"
)

// Service handles image uploading and posting functionality
type Service struct {
    firebaseService *firebase.Service
    httpClient      *http.Client
    createPostURL   string
}

// NewService creates a new upload service instance
func NewService(firebaseService *firebase.Service, httpClient *http.Client, createPostURL string) *Service {
    return &Service{
        firebaseService: firebaseService,
        httpClient:      httpClient,
        createPostURL:   createPostURL,
    }
}

// UploadImage uploads an image to Firebase Storage and creates a post with it
func (s *Service) UploadImage(ctx context.Context, userId, idToken string, imageData []byte, caption string) error {
    log.Printf("Starting image upload for user %s", userId)
    
    // 1. Upload image to Firebase Storage
    imageURL, err := s.firebaseService.UploadFileToStorage(ctx, userId, idToken, imageData, "image/*")
    if err != nil {
        log.Printf("Failed to upload image: %v", err)
        return fmt.Errorf("failed to upload image: %w", err)
    }
    
    log.Printf("Image uploaded successfully, URL: %s", imageURL)
    
    // 2. Create post with the image URL
    postData := map[string]interface{}{
        "data": map[string]interface{}{
            "thumbnail_url": imageURL,
            "caption":      caption,
            "sent_to_all":  true,
        },
    }
    
    jsonData, err := json.Marshal(postData)
    if err != nil {
        return fmt.Errorf("failed to marshal post data: %w", err)
    }
    
    // Create the post request
    req, err := http.NewRequestWithContext(
        ctx, 
        http.MethodPost, 
        s.createPostURL,
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return fmt.Errorf("failed to create post request: %w", err)
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))
    
    log.Printf("Sending request to create post")
    
    // Send the request
    resp, err := s.httpClient.Do(req)
    if err != nil {
        log.Printf("Failed to send post request: %v", err)
        return fmt.Errorf("failed to send post request: %w", err)
    }
    defer resp.Body.Close()
    
    // Check response status
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        log.Printf("Failed to create post: status code %d", resp.StatusCode)
        return fmt.Errorf("failed to create post: status code %d", resp.StatusCode)
    }
    
    log.Printf("Post created successfully")
    return nil
}