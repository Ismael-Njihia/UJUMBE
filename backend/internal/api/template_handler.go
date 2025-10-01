package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/middleware"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TemplateHandler struct {
	db *database.DB
}

func NewTemplateHandler(db *database.DB) *TemplateHandler {
	return &TemplateHandler{db: db}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var template models.EmailTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if template.Name == "" || template.Subject == "" || template.HTMLBody == "" {
		http.Error(w, `{"error":"Name, subject, and HTML body are required"}`, http.StatusBadRequest)
		return
	}

	template.ID = uuid.New()
	template.UserID = userID

	_, err := h.db.Exec(`
		INSERT INTO email_templates (id, user_id, name, subject, html_body, text_body)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, template.ID, template.UserID, template.Name, template.Subject, template.HTMLBody, template.TextBody)
	if err != nil {
		http.Error(w, `{"error":"Failed to create template"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

func (h *TemplateHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.db.Query(`
		SELECT id, user_id, name, subject, html_body, text_body, created_at, updated_at
		FROM email_templates WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch templates"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var templates []models.EmailTemplate
	for rows.Next() {
		var template models.EmailTemplate
		if err := rows.Scan(&template.ID, &template.UserID, &template.Name, &template.Subject,
			&template.HTMLBody, &template.TextBody, &template.CreatedAt, &template.UpdatedAt); err != nil {
			continue
		}
		templates = append(templates, template)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (h *TemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid template ID"}`, http.StatusBadRequest)
		return
	}

	var template models.EmailTemplate
	err = h.db.QueryRow(`
		SELECT id, user_id, name, subject, html_body, text_body, created_at, updated_at
		FROM email_templates WHERE id = $1 AND user_id = $2
	`, templateID, userID).Scan(&template.ID, &template.UserID, &template.Name, &template.Subject,
		&template.HTMLBody, &template.TextBody, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		http.Error(w, `{"error":"Template not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid template ID"}`, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
		DELETE FROM email_templates WHERE id = $1 AND user_id = $2
	`, templateID, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete template"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Template not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Template deleted successfully",
	})
}
