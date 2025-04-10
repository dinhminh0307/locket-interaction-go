package response

// PhoneLoginResponse represents the response from phone login
type PhoneLoginResponse struct {
    Result PhoneLoginResult `json:"result"`
}

// PhoneLoginResult contains the authentication result
type PhoneLoginResult struct {
    Token  string `json:"token"`
    Status int    `json:"status"`
}