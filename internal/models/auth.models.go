package models

import (
    "fmt"
    "strconv"
    "time"
)

// LoginRequest represents the login credentials.
type LoginRequest struct {
    Email            string `json:"email"`
    Password         string `json:"password"`
    ReturnSecureToken bool   `json:"returnSecureToken"`
    ClientType       string `json:"clientType,omitempty"`
}

// FirebaseAuthResponse represents the raw response from Firebase Auth API.
type FirebaseAuthResponse struct {
    Kind           string `json:"kind"`
    LocalID        string `json:"localId"`
    Email          string `json:"email"`
    DisplayName    string `json:"displayName"`
    IDToken        string `json:"idToken"`
    Registered     bool   `json:"registered"`
    ProfilePicture string `json:"profilePicture"`
    RefreshToken   string `json:"refreshToken"`
    ExpiresIn      string `json:"expiresIn"`
}

// ToLoginResponse converts the Firebase response to our standard LoginResponse format
func (f *FirebaseAuthResponse) ToLoginResponse() (*LoginResponse, error) {
    // Convert expiresIn from string to int
    expiresInSeconds, err := strconv.Atoi(f.ExpiresIn)
    if err != nil {
        return nil, fmt.Errorf("invalid expiresIn value: %w", err)
    }

    // Calculate token expiration time
    expiresAt := time.Now().Add(time.Duration(expiresInSeconds) * time.Second)

    return &LoginResponse{
        User: User{
            ID:             f.LocalID,
            Name:           f.DisplayName,
            Email:          f.Email,
            ProfilePicture: f.ProfilePicture,
        },
        Tokens: Tokens{
            Access: Token{
                Token:   f.IDToken,
                Expires: expiresAt,
            },
            Refresh: Token{
                Token:   f.RefreshToken,
                // Refresh tokens typically last longer, but Firebase doesn't provide expiry
                // Setting a default expiration of 30 days for refresh token
                Expires: time.Now().Add(30 * 24 * time.Hour),
            },
        },
        Raw: f, // Store the original response
    }, nil
}

// LoginResponse represents the server response to a successful login,
// reformatted to our application's standard structure.
type LoginResponse struct {
    User   User             `json:"user"`
    Tokens Tokens           `json:"tokens"`
    Raw    interface{}      `json:"-"` // Store the original response but don't include in JSON output
}

// User represents basic user information.
type User struct {
    ID             string `json:"id"`
    Name           string `json:"name"`
    Email          string `json:"email"`
    ProfilePicture string `json:"profilePicture,omitempty"`
}

// Tokens contains both access and refresh tokens.
type Tokens struct {
    Access  Token `json:"access"`
    Refresh Token `json:"refresh"`
}

// Token represents an authentication token.
type Token struct {
    Token   string    `json:"token"`
    Expires time.Time `json:"expires"`
}

// ErrorResponse represents an error from the API.
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
    if e.Details != "" {
        return fmt.Sprintf("API error %d: %s - %s", e.Code, e.Message, e.Details)
    }
    return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// FirebaseErrorResponse represents Firebase-specific error structures
type FirebaseErrorResponse struct {
    Error struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
        Errors  []struct {
            Message string `json:"message"`
            Domain  string `json:"domain"`
            Reason  string `json:"reason"`
        } `json:"errors"`
    } `json:"error"`
}

// ToErrorResponse converts a Firebase error to our standard error format
func (f *FirebaseErrorResponse) ToErrorResponse() *ErrorResponse {
    details := ""
    if len(f.Error.Errors) > 0 {
        details = f.Error.Errors[0].Reason
    }
    
    return &ErrorResponse{
        Code:    f.Error.Code,
        Message: f.Error.Message,
        Details: details,
    }
}