package global

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Initialize loads environment variables from .env file
func Initialize() {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found or could not be loaded: %v", err)
    }

    // log.Printf("FIREBASE_API_KEY: %s", os.Getenv("FIREBASE_API_KEY"))
    // log.Printf("LOCKET_API_KEY: %s", os.Getenv("LOCKET_API_KEY"))
}

// Constants for Locket API
const (
    // LoginURL is the Firebase Authentication endpoint for email/password login
    LoginURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword"
    
    // APIKeyQueryParam is the name of the query parameter for the Firebase API key
    APIKeyQueryParam = "key"

    LoginWithPhone = "https://api.locketcamera.com/signInWithPhonePassword"

    IOSVersion = "1.119.0.1"
)

// GetFirebaseAPIKey returns the Firebase API key from environment variables
func GetFirebaseAPIKey() string {
    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        log.Println("Warning: API_KEY environment variable not set")
    }
    return apiKey
}

// LoginHeaders contains all required HTTP headers for Locket API authentication
var LoginHeaders = map[string]string{
    "Accept":             "*/*",
    "Accept-Encoding":    "gzip, deflate, br",
    "Accept-Language":    "en",
    "Connection":         "keep-alive",
    "Content-Type":       "application/json",
    "Host":               "www.googleapis.com",
    "User-Agent":         "FirebaseAuth.iOS/10.23.1 com.locket.Locket/1.82.0 iPhone/18.0 hw/iPhone12_1",
    "X-Client-Version":   "iOS/FirebaseSDK/10.23.1/FirebaseCore-iOS",
    "X-Firebase-GMPID":   "1:641029076083:ios:cc8eb46290d69b234fa606",
    "X-Ios-Bundle-Identifier": "com.locket.Locket",
}

// FirebaseAuthHeaders contains specialized authentication headers
var FirebaseAuthHeaders = map[string]string{
    "baggage": "sentry-environment=production,sentry-public_key=78fa64317f434fd89d9cc728dd168f50,sentry-release=com.locket.Locket@1.82.0+3,sentry-trace_id=90310ccc8ddd4d059b83321054b6245b",
    "sentry-trace": "90310ccc8ddd4d059b83321054b6245b-3a4920b34e94401d-0",
    "X-Firebase-AppCheck": "eyJraWQiOiJNbjVDS1EiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxOjY0MTAyOTA3NjA4Mzppb3M6Y2M4ZWI0NjI5MGQ2OWIyMzRmYTYwNiIsImF1ZCI6WyJwcm9qZWN0c1wvNjQxMDI5MDc2MDgzIiwicHJvamVjdHNcL2xvY2tldC00MjUyYSJdLCJwcm92aWRlciI6ImRldmljZV9jaGVja19kZXZpY2VfaWRlbnRpZmljYXRpb24iLCJpc3MiOiJodHRwczpcL1wvZmlyZWJhc2VhcHBjaGVjay5nb29nbGVhcGlzLmNvbVwvNjQxMDI5MDc2MDgzIiwiZXhwIjoxNzIyMTY3ODk4LCJpYXQiOjE3MjIxNjQyOTgsImp0aSI6ImlHUGlsT1dDZGg4Mll3UTJXRC1neEpXeWY5TU9RRFhHcU5OR3AzTjFmRGcifQ.lqTOJfdoYLpZwYeeXtRliCdkVT7HMd7_Lj-d44BNTGuxSYPIa9yVAR4upu3vbZSh9mVHYS8kJGYtMqjP-L6YXsk_qsV_gzKC2IhVAV6KbPDRHdevMfBC6fRiOSVn7vt749GVFdZqAuDCXhCILsaMhvgDBgZoDilgAPtpNwyjz-VtRB7OdOUbuKTCqdoSOX0SJWVUMyuI8nH0-unY--YRctunK8JHZDxBaM_ahVggYPWBCpzxq9Yeq8VSPhadG_tGNaADStYPaeeUkZ7DajwWqH5ze6ESpuFNgAigwPxCM735_ZiPeD7zHYwppQA9uqTWszK9v9OvWtFCsgCEe22O8awbNbuEBTKJpDQ8xvZe8iEYyhfUPncER3S-b1CmuXR7tFCdTgQe5j7NGWjFvN_CnL7D2nudLwxWlpqwASCHvHyi8HBaJ5GpgriTLXAAinY48RukRDBi9HwEzpRecELX05KTD2lTOfQCjKyGpfG2VUHP5Xm36YbA3iqTDoDXWMvV",
}

var UPLOADER_HEADERS = map[string]string{
    "content-type": "application/octet-stream",
    "x-goog-upload-protocol": "resumable",
    "x-goog-upload-offset": "0",
    "x-goog-upload-command": "upload, finalize",
    "upload-incomplete": "?0",
    "upload-draft-interop-version": "3",
    "user-agent":
        "com.locket.Locket/1.43.1 iPhone/17.3 hw/iPhone15_3 (GTMSUF/1)",
}