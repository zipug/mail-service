package models

type OTPMessage struct {
	Type      string     `json:"type"`
	Payload   OTPPayload `json:"payload"`
	Timestamp int64      `json:"timestamp"`
}

type OTPMessageType string

const (
	Verify OTPMessageType = "verify"
	Login  OTPMessageType = "login"
)

type OTPPayload struct {
	Type     OTPMessageType `json:"type"`
	UserName string         `json:"username"`
	UserID   int64          `json:"user_id,omitempty"`
	Email    string         `json:"email,omitempty"`
	Code     string         `json:"code,omitempty"`
}
