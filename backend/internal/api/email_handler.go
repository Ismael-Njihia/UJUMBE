package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/middleware"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{emailService: emailService}
}

func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req models.SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.To == "" || req.From == "" {
		http.Error(w, `{"error":"From and To fields are required"}`, http.StatusBadRequest)
		return
	}

	if req.TemplateID == nil && req.Subject == "" {
		http.Error(w, `{"error":"Subject is required when not using a template"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.emailService.SendEmail(userID, req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EmailHandler) GetEmailStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	emailID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid email ID"}`, http.StatusBadRequest)
		return
	}

	var email models.Email
	err = h.emailService.DB.QueryRow(`
		SELECT id, user_id, from_email, to_email, subject, status, ses_message_id, error_message, sent_at, created_at
		FROM emails WHERE id = $1 AND user_id = $2
	`, emailID, userID).Scan(
		&email.ID, &email.UserID, &email.FromEmail, &email.ToEmail, &email.Subject,
		&email.Status, &email.SESMessageID, &email.ErrorMessage, &email.SentAt, &email.CreatedAt,
	)
	if err != nil {
		http.Error(w, `{"error":"Email not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(email)
}

func (h *EmailHandler) GetEmailLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	emailID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid email ID"}`, http.StatusBadRequest)
		return
	}

	// Verify email belongs to user
	var count int
	err = h.emailService.DB.QueryRow(`
		SELECT COUNT(*) FROM emails WHERE id = $1 AND user_id = $2
	`, emailID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, `{"error":"Email not found"}`, http.StatusNotFound)
		return
	}

	rows, err := h.emailService.DB.Query(`
		SELECT id, email_id, event_type, created_at
		FROM email_logs WHERE email_id = $1
		ORDER BY created_at DESC
	`, emailID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch logs"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []models.EmailLog
	for rows.Next() {
		var log models.EmailLog
		if err := rows.Scan(&log.ID, &log.EmailID, &log.EventType, &log.CreatedAt); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *EmailHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	analytics, err := h.emailService.GetAnalytics(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (h *EmailHandler) GetQuota(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	quota, err := h.emailService.GetUserQuota(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quota)
}
