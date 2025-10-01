package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/kafka"
	"github.com/gorilla/mux"
)

// Auth handlers
func (api *API) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := api.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, user, err := api.authService.Login(req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (api *API) getUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	user, err := api.authService.GetUser(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (api *API) getDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	stats, err := api.emailService.GetDashboardStats(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// API Key handlers
func (api *API) getAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	apiKeys, err := api.authService.GetUserAPIKeys(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get API keys")
		return
	}

	respondJSON(w, http.StatusOK, apiKeys)
}

func (api *API) createAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	apiKey, err := api.authService.CreateAPIKey(userID, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	respondJSON(w, http.StatusCreated, apiKey)
}

func (api *API) revokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	keyID := vars["id"]

	if err := api.authService.RevokeAPIKey(keyID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to revoke API key")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "API key revoked successfully"})
}

// Domain handlers
func (api *API) getDomains(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	domains, err := api.domainService.GetUserDomains(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get domains")
		return
	}

	respondJSON(w, http.StatusOK, domains)
}

func (api *API) addDomain(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req struct {
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	domain, err := api.domainService.AddDomain(userID, req.Domain)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, domain)
}

func (api *API) getDomain(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	domainID := vars["id"]

	domain, err := api.domainService.GetDomain(domainID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	respondJSON(w, http.StatusOK, domain)
}

func (api *API) verifyDomain(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	domainID := vars["id"]

	if err := api.domainService.VerifyDomain(domainID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to verify domain")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Domain verified successfully"})
}

func (api *API) deleteDomain(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	domainID := vars["id"]

	if err := api.domainService.DeleteDomain(domainID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete domain")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Domain deleted successfully"})
}

// Template handlers
func (api *API) getTemplates(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	templates, err := api.emailService.GetUserTemplates(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get templates")
		return
	}

	respondJSON(w, http.StatusOK, templates)
}

func (api *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req struct {
		Name      string            `json:"name"`
		Subject   string            `json:"subject"`
		HTMLBody  string            `json:"html_body"`
		TextBody  string            `json:"text_body"`
		Variables map[string]string `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	template, err := api.emailService.CreateTemplate(userID, req.Name, req.Subject, req.HTMLBody, req.TextBody, req.Variables)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create template")
		return
	}

	respondJSON(w, http.StatusCreated, template)
}

func (api *API) getTemplate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	templateID := vars["id"]

	template, err := api.emailService.GetTemplate(templateID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Template not found")
		return
	}

	respondJSON(w, http.StatusOK, template)
}

func (api *API) updateTemplate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	templateID := vars["id"]

	var req struct {
		Name      string            `json:"name"`
		Subject   string            `json:"subject"`
		HTMLBody  string            `json:"html_body"`
		TextBody  string            `json:"text_body"`
		Variables map[string]string `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := api.emailService.UpdateTemplate(templateID, userID, req.Name, req.Subject, req.HTMLBody, req.TextBody, req.Variables); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update template")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Template updated successfully"})
}

func (api *API) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	vars := mux.Vars(r)
	templateID := vars["id"]

	if err := api.emailService.DeleteTemplate(templateID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete template")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}

// Email handlers
func (api *API) sendEmail(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req struct {
		TemplateID string            `json:"template_id"`
		From       string            `json:"from"`
		To         string            `json:"to"`
		Subject    string            `json:"subject"`
		HTMLBody   string            `json:"html_body"`
		TextBody   string            `json:"text_body"`
		Variables  map[string]string `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check quota
	hasQuota, err := api.emailService.CheckQuota(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to check quota")
		return
	}
	if !hasQuota {
		respondError(w, http.StatusForbidden, "Insufficient email quota")
		return
	}

	subject := req.Subject
	htmlBody := req.HTMLBody
	textBody := req.TextBody

	// If template ID is provided, use template
	if req.TemplateID != "" {
		template, err := api.emailService.GetTemplate(req.TemplateID, userID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Template not found")
			return
		}
		subject, htmlBody, textBody = api.emailService.RenderTemplate(template, req.Variables)
	}

	// Create email log
	logID := ""
	if err := api.emailService.LogEmail(userID, req.TemplateID, req.From, req.To, subject, "queued", "", ""); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to log email")
		return
	}

	// Queue email in Kafka
	job := &kafka.EmailJob{
		ID:         logID,
		UserID:     userID,
		TemplateID: req.TemplateID,
		From:       req.From,
		To:         req.To,
		Subject:    subject,
		HTMLBody:   htmlBody,
		TextBody:   textBody,
		Variables:  req.Variables,
	}

	if err := api.kafkaProducer.PublishEmail(r.Context(), job); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to queue email")
		return
	}

	// Increment emails sent counter
	if err := api.emailService.IncrementEmailsSent(userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update quota")
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "Email queued successfully",
		"status":  "queued",
	})
}

func (api *API) getEmailLogs(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	logs, err := api.emailService.GetEmailLogs(userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get email logs")
		return
	}

	respondJSON(w, http.StatusOK, logs)
}

// Analytics handlers
func (api *API) getAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	analytics, err := api.emailService.GetAnalytics(userID, startDate, endDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get analytics")
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}

// Billing handlers
func (api *API) getTransactions(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	transactions, err := api.billingService.GetUserTransactions(userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

func (api *API) initiateTopup(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req struct {
		Amount      float64 `json:"amount"`
		PhoneNumber string  `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Calculate emails to purchase
	emailsPurchased := int(req.Amount / api.config.EmailPricePerUnit)

	// Create transaction
	transaction, err := api.billingService.CreateTransaction(userID, req.Amount, "mpesa", req.PhoneNumber, emailsPurchased)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Initiate M-Pesa STK push
	stkResponse, err := api.mpesaClient.InitiateSTKPush(
		req.PhoneNumber,
		req.Amount,
		transaction.ID,
		"Email credits top-up",
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to initiate payment")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"transaction_id": transaction.ID,
		"mpesa_response": stkResponse,
		"message":        "Payment initiated. Please complete on your phone",
	})
}

func (api *API) mpesaCallback(w http.ResponseWriter, r *http.Request) {
	var callbackData json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&callbackData); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid callback data")
		return
	}

	callback, err := api.mpesaClient.ParseCallback(callbackData)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse callback")
		return
	}

	// Extract transaction details
	merchantRequestID := callback.Body.StkCallback.MerchantRequestID
	resultCode := callback.Body.StkCallback.ResultCode

	if resultCode == 0 {
		// Payment successful
		var mpesaReceiptNo string
		for _, item := range callback.Body.StkCallback.CallbackMetadata.Item {
			if item.Name == "MpesaReceiptNumber" {
				mpesaReceiptNo = item.Value.(string)
				break
			}
		}

		// Update transaction status
		transaction, err := api.billingService.GetTransaction(merchantRequestID)
		if err == nil {
			api.billingService.UpdateTransactionStatus(transaction.ID, "completed", mpesaReceiptNo)
			api.billingService.AddEmailQuota(transaction.UserID, transaction.EmailsPurchased)
		}
	} else {
		// Payment failed
		transaction, err := api.billingService.GetTransaction(merchantRequestID)
		if err == nil {
			api.billingService.UpdateTransactionStatus(transaction.ID, "failed", "")
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Callback processed"})
}
