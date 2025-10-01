package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/kafka"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/ses"
	"github.com/Ismael-Njihia/UJUMBE/backend/pkg/utils"
	"github.com/google/uuid"
)

type EmailService struct {
	DB       *database.DB
	ses      *ses.SESClient
	producer *kafka.Producer
}

func NewEmailService(db *database.DB, sesClient *ses.SESClient, producer *kafka.Producer) *EmailService {
	return &EmailService{
		DB:       db,
		ses:      sesClient,
		producer: producer,
	}
}

func (s *EmailService) SendEmail(userID uuid.UUID, req models.SendEmailRequest) (*models.SendEmailResponse, error) {
	// Check and deduct quota
	remaining, err := s.DeductQuota(userID)
	if err != nil {
		return nil, err
	}

	// Prepare email data
	var htmlBody, textBody, subject string
	var templateID *uuid.UUID

	if req.TemplateID != nil {
		// Use template
		template, err := s.GetTemplate(userID, *req.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("template not found: %w", err)
		}

		subject = template.Subject
		htmlBody = utils.ReplaceTemplateVariables(template.HTMLBody, req.TemplateData)
		textBody = utils.ReplaceTemplateVariables(template.TextBody, req.TemplateData)
		templateID = req.TemplateID
	} else {
		// Direct email
		subject = req.Subject
		htmlBody = req.HTMLBody
		textBody = req.TextBody
	}

	// Validate sender domain
	fromDomain := strings.Split(req.From, "@")[1]
	if !s.IsDomainVerified(userID, fromDomain) {
		return nil, fmt.Errorf("sender domain not verified: %s", fromDomain)
	}

	// Create email record
	emailID := uuid.New()
	_, err = s.DB.Exec(`
		INSERT INTO emails (id, user_id, template_id, from_email, to_email, subject, html_body, text_body, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, emailID, userID, templateID, req.From, req.To, subject, htmlBody, textBody, "pending")
	if err != nil {
		return nil, fmt.Errorf("failed to create email record: %w", err)
	}

	// Create log entry
	s.CreateEmailLog(emailID, "created", nil)

	// Send to Kafka for processing
	job := kafka.EmailJob{
		EmailID:  emailID,
		UserID:   userID,
		From:     req.From,
		To:       req.To,
		Subject:  subject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}

	err = s.producer.ProduceEmailJob(job)
	if err != nil {
		// Update status to failed
		s.DB.Exec("UPDATE emails SET status = $1, error_message = $2 WHERE id = $3", "failed", err.Error(), emailID)
		return nil, fmt.Errorf("failed to queue email: %w", err)
	}

	s.CreateEmailLog(emailID, "queued", nil)

	return &models.SendEmailResponse{
		Success:   true,
		EmailID:   emailID,
		Message:   "Email queued successfully",
		Remaining: remaining,
	}, nil
}

func (s *EmailService) ProcessEmailJob(job kafka.EmailJob) error {
	// Send via SES
	messageID, err := s.ses.SendEmail(job.From, job.To, job.Subject, job.HTMLBody, job.TextBody)
	if err != nil {
		// Update status to failed
		s.DB.Exec("UPDATE emails SET status = $1, error_message = $2 WHERE id = $3", "failed", err.Error(), job.EmailID)
		s.CreateEmailLog(job.EmailID, "failed", map[string]interface{}{"error": err.Error()})
		return err
	}

	// Update status to sent
	now := time.Now()
	_, err = s.DB.Exec("UPDATE emails SET status = $1, ses_message_id = $2, sent_at = $3 WHERE id = $4", "sent", messageID, now, job.EmailID)
	if err != nil {
		return fmt.Errorf("failed to update email status: %w", err)
	}

	s.CreateEmailLog(job.EmailID, "sent", map[string]interface{}{"message_id": messageID})
	return nil
}

func (s *EmailService) DeductQuota(userID uuid.UUID) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Get current quota
	var quota models.UserQuota
	err = tx.QueryRow(`
		SELECT id, user_id, free_emails_remaining, paid_emails_balance, monthly_reset_date
		FROM user_quotas WHERE user_id = $1 FOR UPDATE
	`, userID).Scan(&quota.ID, &quota.UserID, &quota.FreeEmailsRemaining, &quota.PaidEmailsBalance, &quota.MonthlyResetDate)

	if err == sql.ErrNoRows {
		// Create new quota
		_, err = tx.Exec(`
			INSERT INTO user_quotas (user_id, free_emails_remaining, paid_emails_balance, monthly_reset_date)
			VALUES ($1, 100, 0, CURRENT_DATE)
		`, userID)
		if err != nil {
			return 0, err
		}
		quota.FreeEmailsRemaining = 100
		quota.PaidEmailsBalance = 0
	} else if err != nil {
		return 0, err
	}

	// Check if monthly reset is needed
	if time.Now().After(quota.MonthlyResetDate.AddDate(0, 1, 0)) {
		quota.FreeEmailsRemaining = 100
		quota.MonthlyResetDate = time.Now()
	}

	// Deduct from quota
	if quota.FreeEmailsRemaining > 0 {
		quota.FreeEmailsRemaining--
	} else if quota.PaidEmailsBalance > 0 {
		quota.PaidEmailsBalance--
	} else {
		return 0, fmt.Errorf("insufficient quota")
	}

	// Update quota
	_, err = tx.Exec(`
		UPDATE user_quotas 
		SET free_emails_remaining = $1, paid_emails_balance = $2, monthly_reset_date = $3
		WHERE user_id = $4
	`, quota.FreeEmailsRemaining, quota.PaidEmailsBalance, quota.MonthlyResetDate, userID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	remaining := quota.FreeEmailsRemaining + quota.PaidEmailsBalance
	return remaining, nil
}

func (s *EmailService) IsDomainVerified(userID uuid.UUID, domain string) bool {
	var verified bool
	err := s.DB.QueryRow(`
		SELECT verified FROM verified_domains 
		WHERE user_id = $1 AND domain = $2 AND verified = true
	`, userID, domain).Scan(&verified)
	return err == nil && verified
}

func (s *EmailService) GetTemplate(userID uuid.UUID, templateID uuid.UUID) (*models.EmailTemplate, error) {
	var template models.EmailTemplate
	err := s.DB.QueryRow(`
		SELECT id, user_id, name, subject, html_body, text_body, created_at, updated_at
		FROM email_templates WHERE id = $1 AND user_id = $2
	`, templateID, userID).Scan(
		&template.ID, &template.UserID, &template.Name, &template.Subject,
		&template.HTMLBody, &template.TextBody, &template.CreatedAt, &template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (s *EmailService) CreateEmailLog(emailID uuid.UUID, eventType string, eventData map[string]interface{}) error {
	// For simplicity, we'll store event data as a string representation
	// In production, you'd use JSONB
	_, err := s.DB.Exec(`
		INSERT INTO email_logs (id, email_id, event_type, event_data)
		VALUES ($1, $2, $3, NULL)
	`, uuid.New(), emailID, eventType)
	return err
}

func (s *EmailService) GetUserQuota(userID uuid.UUID) (*models.UserQuota, error) {
	var quota models.UserQuota
	err := s.DB.QueryRow(`
		SELECT id, user_id, free_emails_remaining, paid_emails_balance, monthly_reset_date, created_at, updated_at
		FROM user_quotas WHERE user_id = $1
	`, userID).Scan(
		&quota.ID, &quota.UserID, &quota.FreeEmailsRemaining, &quota.PaidEmailsBalance,
		&quota.MonthlyResetDate, &quota.CreatedAt, &quota.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		// Create default quota
		_, err = s.DB.Exec(`
			INSERT INTO user_quotas (user_id, free_emails_remaining, paid_emails_balance, monthly_reset_date)
			VALUES ($1, 100, 0, CURRENT_DATE)
		`, userID)
		if err != nil {
			return nil, err
		}
		return s.GetUserQuota(userID)
	}
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

func (s *EmailService) GetAnalytics(userID uuid.UUID) (*models.AnalyticsResponse, error) {
	var analytics models.AnalyticsResponse

	// Get email counts
	err := s.DB.QueryRow(`
		SELECT 
			COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending
		FROM emails WHERE user_id = $1
	`, userID).Scan(&analytics.TotalEmailsSent, &analytics.TotalEmailsFailed, &analytics.TotalEmailsPending)
	if err != nil {
		return nil, err
	}

	// Get quota
	quota, err := s.GetUserQuota(userID)
	if err != nil {
		return nil, err
	}
	analytics.FreeEmailsRemaining = quota.FreeEmailsRemaining
	analytics.PaidEmailsBalance = quota.PaidEmailsBalance

	// Calculate success rate
	total := analytics.TotalEmailsSent + analytics.TotalEmailsFailed
	if total > 0 {
		analytics.SuccessRate = float64(analytics.TotalEmailsSent) / float64(total) * 100
	}

	return &analytics, nil
}
