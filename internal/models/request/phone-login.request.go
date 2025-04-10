package request

// PhoneLoginRequest represents the request body for phone number login
type PhoneLoginRequest struct {
    Data PhoneLoginData `json:"data"`
}

// PhoneLoginData contains the phone login credentials
type PhoneLoginData struct {
    Phone      string `json:"phone"`
    Password   string `json:"password"`
    IOSVersion string `json:"ios_version"`
}