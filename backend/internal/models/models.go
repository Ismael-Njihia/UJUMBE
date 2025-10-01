package models

import (
	"time"
)

// User represents a registered user
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	EmailQuota   int       `json:"email_quota"`
	EmailsSent   int       `json:"emails_sent"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Domain represents a verified sender domain
type Domain struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Domain           string    `json:"domain"`
	IsVerified       bool      `json:"is_verified"`
	VerificationCode string    `json:"verification_code,omitempty"`
	DKIMStatus       string    `json:"dkim_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Subject   string            `json:"subject"`
	HTMLBody  string            `json:"html_body"`
	TextBody  string            `json:"text_body"`
	Variables map[string]string `json:"variables,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// EmailLog represents a sent email log
type EmailLog struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TemplateID   string    `json:"template_id,omitempty"`
	FromEmail    string    `json:"from_email"`
	ToEmail      string    `json:"to_email"`
	Subject      string    `json:"subject"`
	Status       string    `json:"status"` // queued, sent, failed, bounced
	ErrorMessage string    `json:"error_message,omitempty"`
	MessageID    string    `json:"message_id,omitempty"`
	SentAt       time.Time `json:"sent_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Transaction represents a billing transaction
type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Type            string    `json:"type"`   // mpesa, credit
	Status          string    `json:"status"` // pending, completed, failed
	MpesaReceiptNo  string    `json:"mpesa_receipt_no,omitempty"`
	PhoneNumber     string    `json:"phone_number,omitempty"`
	EmailsPurchased int       `json:"emails_purchased"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Analytics represents email analytics
type Analytics struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Date       time.Time `json:"date"`
	EmailsSent int       `json:"emails_sent"`
	Failed     int       `json:"failed"`
	Bounced    int       `json:"bounced"`
	CreatedAt  time.Time `json:"created_at"`
}
