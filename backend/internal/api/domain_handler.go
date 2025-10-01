package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/middleware"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/ses"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DomainHandler struct {
	db  *database.DB
	ses *ses.SESClient
}

func NewDomainHandler(db *database.DB, sesClient *ses.SESClient) *DomainHandler {
	return &DomainHandler{db: db, ses: sesClient}
}

func (h *DomainHandler) AddDomain(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, `{"error":"Domain is required"}`, http.StatusBadRequest)
		return
	}

	// Generate verification token
	verificationToken := uuid.New().String()

	// Add domain to database
	domainID := uuid.New()
	_, err := h.db.Exec(`
		INSERT INTO verified_domains (id, user_id, domain, verified, verification_token)
		VALUES ($1, $2, $3, $4, $5)
	`, domainID, userID, req.Domain, false, verificationToken)
	if err != nil {
		http.Error(w, `{"error":"Failed to add domain"}`, http.StatusInternalServerError)
		return
	}

	// Initiate SES domain verification
	err = h.ses.VerifyDomain(req.Domain)
	if err != nil {
		http.Error(w, `{"error":"Failed to initiate domain verification with SES"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":             true,
		"message":             "Domain added. Please verify ownership by adding DNS records.",
		"domain_id":           domainID,
		"verification_token":  verificationToken,
	})
}

func (h *DomainHandler) GetDomains(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := h.db.Query(`
		SELECT id, user_id, domain, verified, created_at, updated_at
		FROM verified_domains WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch domains"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var domains []models.VerifiedDomain
	for rows.Next() {
		var domain models.VerifiedDomain
		if err := rows.Scan(&domain.ID, &domain.UserID, &domain.Domain, &domain.Verified,
			&domain.CreatedAt, &domain.UpdatedAt); err != nil {
			continue
		}
		domains = append(domains, domain)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (h *DomainHandler) VerifyDomain(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	domainID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid domain ID"}`, http.StatusBadRequest)
		return
	}

	// Update domain as verified
	result, err := h.db.Exec(`
		UPDATE verified_domains SET verified = true
		WHERE id = $1 AND user_id = $2
	`, domainID, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to verify domain"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Domain not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Domain verified successfully",
	})
}

func (h *DomainHandler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	domainID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid domain ID"}`, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
		DELETE FROM verified_domains WHERE id = $1 AND user_id = $2
	`, domainID, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete domain"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Domain not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Domain deleted successfully",
	})
}
