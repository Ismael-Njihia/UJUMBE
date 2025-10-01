package email

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrInsufficientQuota = errors.New("insufficient email quota")
	ErrTemplateNotFound  = errors.New("template not found")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateTemplate creates a new email template
func (s *Service) CreateTemplate(userID, name, subject, htmlBody, textBody string, variables map[string]string) (*models.EmailTemplate, error) {
	template := &models.EmailTemplate{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Subject:   subject,
		HTMLBody:  htmlBody,
		TextBody:  textBody,
		Variables: variables,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(
		`INSERT INTO email_templates (id, user_id, name, subject, html_body, text_body, variables, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		template.ID, template.UserID, template.Name, template.Subject, template.HTMLBody, template.TextBody, variablesJSON, template.CreatedAt, template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (s *Service) GetTemplate(templateID, userID string) (*models.EmailTemplate, error) {
	var template models.EmailTemplate
	var variablesJSON []byte

	err := s.db.QueryRow(
		`SELECT id, user_id, name, subject, html_body, text_body, variables, created_at, updated_at
		 FROM email_templates WHERE id = $1 AND user_id = $2`,
		templateID, userID,
	).Scan(&template.ID, &template.UserID, &template.Name, &template.Subject, &template.HTMLBody, &template.TextBody, &variablesJSON, &template.CreatedAt, &template.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrTemplateNotFound
	}
	if err != nil {
		return nil, err
	}

	if len(variablesJSON) > 0 {
		if err := json.Unmarshal(variablesJSON, &template.Variables); err != nil {
			return nil, err
		}
	}

	return &template, nil
}

// GetUserTemplates retrieves all templates for a user
func (s *Service) GetUserTemplates(userID string) ([]models.EmailTemplate, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, name, subject, html_body, text_body, variables, created_at, updated_at
		 FROM email_templates WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.EmailTemplate
	for rows.Next() {
		var template models.EmailTemplate
		var variablesJSON []byte

		if err := rows.Scan(&template.ID, &template.UserID, &template.Name, &template.Subject, &template.HTMLBody, &template.TextBody, &variablesJSON, &template.CreatedAt, &template.UpdatedAt); err != nil {
			return nil, err
		}

		if len(variablesJSON) > 0 {
			if err := json.Unmarshal(variablesJSON, &template.Variables); err != nil {
				return nil, err
			}
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// UpdateTemplate updates an existing template
func (s *Service) UpdateTemplate(templateID, userID, name, subject, htmlBody, textBody string, variables map[string]string) error {
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`UPDATE email_templates SET name = $1, subject = $2, html_body = $3, text_body = $4, variables = $5, updated_at = $6
		 WHERE id = $7 AND user_id = $8`,
		name, subject, htmlBody, textBody, variablesJSON, time.Now(), templateID, userID,
	)
	return err
}

// DeleteTemplate deletes a template
func (s *Service) DeleteTemplate(templateID, userID string) error {
	_, err := s.db.Exec(
		`DELETE FROM email_templates WHERE id = $1 AND user_id = $2`,
		templateID, userID,
	)
	return err
}

// RenderTemplate replaces variables in template with actual values
func (s *Service) RenderTemplate(template *models.EmailTemplate, variables map[string]string) (string, string, string) {
	subject := template.Subject
	htmlBody := template.HTMLBody
	textBody := template.TextBody

	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		subject = strings.ReplaceAll(subject, placeholder, value)
		htmlBody = strings.ReplaceAll(htmlBody, placeholder, value)
		textBody = strings.ReplaceAll(textBody, placeholder, value)
	}

	return subject, htmlBody, textBody
}

// LogEmail creates an email log entry
func (s *Service) LogEmail(userID, templateID, fromEmail, toEmail, subject, status, errorMessage, messageID string) error {
	logID := uuid.New().String()
	now := time.Now()

	var sentAt *time.Time
	if status == "sent" {
		sentAt = &now
	}

	_, err := s.db.Exec(
		`INSERT INTO email_logs (id, user_id, template_id, from_email, to_email, subject, status, error_message, message_id, sent_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		logID, userID, templateID, fromEmail, toEmail, subject, status, errorMessage, messageID, sentAt, now,
	)
	return err
}

// GetEmailLogs retrieves email logs for a user
func (s *Service) GetEmailLogs(userID string, limit, offset int) ([]models.EmailLog, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, template_id, from_email, to_email, subject, status, error_message, message_id, sent_at, created_at
		 FROM email_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.EmailLog
	for rows.Next() {
		var log models.EmailLog
		var templateID, errorMessage, messageID sql.NullString
		var sentAt sql.NullTime

		if err := rows.Scan(&log.ID, &log.UserID, &templateID, &log.FromEmail, &log.ToEmail, &log.Subject, &log.Status, &errorMessage, &messageID, &sentAt, &log.CreatedAt); err != nil {
			return nil, err
		}

		if templateID.Valid {
			log.TemplateID = templateID.String
		}
		if errorMessage.Valid {
			log.ErrorMessage = errorMessage.String
		}
		if messageID.Valid {
			log.MessageID = messageID.String
		}
		if sentAt.Valid {
			log.SentAt = sentAt.Time
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// CheckQuota verifies if user has sufficient email quota
func (s *Service) CheckQuota(userID string) (bool, error) {
	var emailQuota, emailsSent int
	err := s.db.QueryRow(
		`SELECT email_quota, emails_sent FROM users WHERE id = $1`,
		userID,
	).Scan(&emailQuota, &emailsSent)

	if err != nil {
		return false, err
	}

	return emailsSent < emailQuota, nil
}

// IncrementEmailsSent increments the emails_sent counter for a user
func (s *Service) IncrementEmailsSent(userID string) error {
	_, err := s.db.Exec(
		`UPDATE users SET emails_sent = emails_sent + 1, updated_at = $1 WHERE id = $2`,
		time.Now(), userID,
	)
	return err
}

// UpdateEmailLogStatus updates the status of an email log
func (s *Service) UpdateEmailLogStatus(logID, status, errorMessage, messageID string) error {
	now := time.Now()
	var sentAt *time.Time
	if status == "sent" {
		sentAt = &now
	}

	_, err := s.db.Exec(
		`UPDATE email_logs SET status = $1, error_message = $2, message_id = $3, sent_at = $4 WHERE id = $5`,
		status, errorMessage, messageID, sentAt, logID,
	)
	return err
}

// RecordAnalytics records daily analytics
func (s *Service) RecordAnalytics(userID string, emailsSent, failed, bounced int) error {
	date := time.Now().Truncate(24 * time.Hour)
	analyticsID := uuid.New().String()

	_, err := s.db.Exec(
		`INSERT INTO analytics (id, user_id, date, emails_sent, failed, bounced, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id, date) DO UPDATE SET
		 emails_sent = analytics.emails_sent + $4,
		 failed = analytics.failed + $5,
		 bounced = analytics.bounced + $6`,
		analyticsID, userID, date, emailsSent, failed, bounced, time.Now(),
	)
	return err
}

// GetAnalytics retrieves analytics for a user
func (s *Service) GetAnalytics(userID string, startDate, endDate time.Time) ([]models.Analytics, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, date, emails_sent, failed, bounced, created_at
		 FROM analytics WHERE user_id = $1 AND date BETWEEN $2 AND $3 ORDER BY date DESC`,
		userID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []models.Analytics
	for rows.Next() {
		var record models.Analytics
		if err := rows.Scan(&record.ID, &record.UserID, &record.Date, &record.EmailsSent, &record.Failed, &record.Bounced, &record.CreatedAt); err != nil {
			return nil, err
		}
		analytics = append(analytics, record)
	}

	return analytics, nil
}

// GetDashboardStats retrieves overall statistics for dashboard
func (s *Service) GetDashboardStats(userID string) (map[string]interface{}, error) {
	var totalSent, totalFailed, totalBounced int
	var emailQuota, emailsSent int
	var balance float64

	// Get user stats
	err := s.db.QueryRow(
		`SELECT email_quota, emails_sent, balance FROM users WHERE id = $1`,
		userID,
	).Scan(&emailQuota, &emailsSent, &balance)
	if err != nil {
		return nil, err
	}

	// Get total counts from logs
	err = s.db.QueryRow(
		`SELECT 
			COUNT(CASE WHEN status = 'sent' THEN 1 END) as total_sent,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as total_failed,
			COUNT(CASE WHEN status = 'bounced' THEN 1 END) as total_bounced
		 FROM email_logs WHERE user_id = $1`,
		userID,
	).Scan(&totalSent, &totalFailed, &totalBounced)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"email_quota":     emailQuota,
		"emails_sent":     emailsSent,
		"balance":         balance,
		"total_sent":      totalSent,
		"total_failed":    totalFailed,
		"total_bounced":   totalBounced,
		"quota_remaining": emailQuota - emailsSent,
	}

	return stats, nil
}
