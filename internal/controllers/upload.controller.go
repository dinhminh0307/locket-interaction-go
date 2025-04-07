package controllers

import (
    "encoding/json"
    "io"
    "log"
    "net/http"

    "locket-interaction-go/internal/services/upload"
    "locket-interaction-go/internal/middlewares"
)

// UploadController handles HTTP requests related to image uploads
type UploadController struct {
    uploadService *upload.Service
}

// NewUploadController creates a new upload controller
func NewUploadController(uploadService *upload.Service) *UploadController {
    return &UploadController{
        uploadService: uploadService,
    }
}

// HandleUploadImage returns a handler function for the image upload endpoint
func (c *UploadController) HandleUploadImage() http.HandlerFunc {
    return middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Get user information from context (set by auth middleware)
        userID := r.Context().Value("userId").(string)
        idToken := r.Context().Value("idToken").(string)

		log.Printf("Adding user ID to context: %s", userID)
		log.Printf("Adding auth token to context: %s", idToken)

        // Parse multipart form data (max 10MB)
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            log.Printf("Error parsing multipart form: %v", err)
            http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Get image file
        file, _, err := r.FormFile("image")
        if err != nil {
            log.Printf("Error getting image file: %v", err)
            http.Error(w, "Bad request: image field is required", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Read image data
        imageData, err := io.ReadAll(file)
        if err != nil {
            log.Printf("Error reading image data: %v", err)
            http.Error(w, "Failed to read image data", http.StatusInternalServerError)
            return
        }

        // Get caption from form
        caption := r.FormValue("caption")

        // Upload image and create post
        err = c.uploadService.UploadImage(r.Context(), userID, idToken, imageData, caption)
        if err != nil {
            log.Printf("Error uploading image: %v", err)
            http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Return success response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        
        response := map[string]string{"status": "success", "message": "Image uploaded successfully"}
        json.NewEncoder(w).Encode(response)
    })
}