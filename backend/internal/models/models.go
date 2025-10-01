package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	APIKey       string    `json:"api_key" db:"api_key"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserQuota struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	UserID               uuid.UUID `json:"user_id" db:"user_id"`
	FreeEmailsRemaining  int       `json:"free_emails_remaining" db:"free_emails_remaining"`
	PaidEmailsBalance    int       `json:"paid_emails_balance" db:"paid_emails_balance"`
	MonthlyResetDate     time.Time `json:"monthly_reset_date" db:"monthly_reset_date"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type VerifiedDomain struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	Domain            string    `json:"domain" db:"domain"`
	Verified          bool      `json:"verified" db:"verified"`
	VerificationToken string    `json:"verification_token,omitempty" db:"verification_token"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type EmailTemplate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Subject   string    `json:"subject" db:"subject"`
	HTMLBody  string    `json:"html_body" db:"html_body"`
	TextBody  string    `json:"text_body,omitempty" db:"text_body"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Email struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	TemplateID   *uuid.UUID `json:"template_id,omitempty" db:"template_id"`
	FromEmail    string     `json:"from_email" db:"from_email"`
	ToEmail      string     `json:"to_email" db:"to_email"`
	Subject      string     `json:"subject" db:"subject"`
	HTMLBody     string     `json:"html_body,omitempty" db:"html_body"`
	TextBody     string     `json:"text_body,omitempty" db:"text_body"`
	Status       string     `json:"status" db:"status"`
	SESMessageID string     `json:"ses_message_id,omitempty" db:"ses_message_id"`
	ErrorMessage string     `json:"error_message,omitempty" db:"error_message"`
	SentAt       *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type EmailLog struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	EmailID   uuid.UUID              `json:"email_id" db:"email_id"`
	EventType string                 `json:"event_type" db:"event_type"`
	EventData map[string]interface{} `json:"event_data,omitempty" db:"event_data"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

type Transaction struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	UserID              uuid.UUID `json:"user_id" db:"user_id"`
	Amount              float64   `json:"amount" db:"amount"`
	PhoneNumber         string    `json:"phone_number" db:"phone_number"`
	MpesaReceiptNumber  string    `json:"mpesa_receipt_number,omitempty" db:"mpesa_receipt_number"`
	TransactionType     string    `json:"transaction_type" db:"transaction_type"`
	Status              string    `json:"status" db:"status"`
	EmailsCredited      int       `json:"emails_credited" db:"emails_credited"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// API Request/Response models
type SendEmailRequest struct {
	TemplateID     *uuid.UUID             `json:"template_id,omitempty"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	Subject        string                 `json:"subject,omitempty"`
	HTMLBody       string                 `json:"html_body,omitempty"`
	TextBody       string                 `json:"text_body,omitempty"`
	TemplateData   map[string]interface{} `json:"template_data,omitempty"`
}

type SendEmailResponse struct {
	Success   bool      `json:"success"`
	EmailID   uuid.UUID `json:"email_id,omitempty"`
	Message   string    `json:"message"`
	Remaining int       `json:"remaining"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	APIKey  string `json:"api_key,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	APIKey  string `json:"api_key,omitempty"`
	Message string `json:"message"`
}

type MpesaTopupRequest struct {
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
}

type MpesaTopupResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	CheckoutRequestID string `json:"checkout_request_id,omitempty"`
}

type AnalyticsResponse struct {
	TotalEmailsSent     int     `json:"total_emails_sent"`
	TotalEmailsFailed   int     `json:"total_emails_failed"`
	TotalEmailsPending  int     `json:"total_emails_pending"`
	FreeEmailsRemaining int     `json:"free_emails_remaining"`
	PaidEmailsBalance   int     `json:"paid_emails_balance"`
	SuccessRate         float64 `json:"success_rate"`
}
