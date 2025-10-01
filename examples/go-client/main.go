package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// Update this to your UJUMBE API URL
	UJUMBE_API_URL = "http://localhost:8080/api/v1"
	// Add your API key here
	API_KEY = "your-api-key-here"
)

type EmailRequest struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Subject    string                 `json:"subject,omitempty"`
	HTMLBody   string                 `json:"html_body,omitempty"`
	TextBody   string                 `json:"text_body,omitempty"`
	TemplateID *string                `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
}

type EmailResponse struct {
	Success   bool   `json:"success"`
	EmailID   string `json:"email_id,omitempty"`
	Message   string `json:"message"`
	Remaining int    `json:"remaining"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// UjumbeClient is a simple client for the UJUMBE API
type UjumbeClient struct {
	apiURL string
	apiKey string
	client *http.Client
}

// NewUjumbeClient creates a new UJUMBE API client
func NewUjumbeClient(apiURL, apiKey string) *UjumbeClient {
	return &UjumbeClient{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// SendEmail sends an email using UJUMBE
func (c *UjumbeClient) SendEmail(req EmailRequest) (*EmailResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.apiURL+"/emails/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errResp.Error)
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &emailResp, nil
}

func main() {
	// Create UJUMBE client
	client := NewUjumbeClient(UJUMBE_API_URL, API_KEY)

	// Example 1: Send a simple email
	fmt.Println("=== Example 1: Simple Email ===")
	simpleEmail := EmailRequest{
		From:     "hello@yourdomain.com",
		To:       "user@example.com",
		Subject:  "Welcome to Our Service",
		HTMLBody: "<h1>Welcome!</h1><p>Thank you for signing up.</p>",
		TextBody: "Welcome! Thank you for signing up.",
	}

	resp, err := client.SendEmail(simpleEmail)
	if err != nil {
		fmt.Printf("Error sending simple email: %v\n", err)
	} else {
		fmt.Printf("Success! Email ID: %s, Remaining quota: %d\n", resp.EmailID, resp.Remaining)
	}

	// Example 2: Send email with template (requires template ID)
	// First create a template through the dashboard or API
	fmt.Println("\n=== Example 2: Template Email ===")
	templateID := "your-template-uuid-here" // Replace with actual template ID
	templateEmail := EmailRequest{
		From:       "hello@yourdomain.com",
		To:         "user@example.com",
		TemplateID: &templateID,
		TemplateData: map[string]interface{}{
			"name":     "John Doe",
			"company":  "ACME Inc",
			"action_url": "https://example.com/verify",
		},
	}

	resp, err = client.SendEmail(templateEmail)
	if err != nil {
		fmt.Printf("Error sending template email: %v\n", err)
	} else {
		fmt.Printf("Success! Email ID: %s, Remaining quota: %d\n", resp.EmailID, resp.Remaining)
	}

	// Example 3: Send transactional email
	fmt.Println("\n=== Example 3: Transactional Email ===")
	transactionalEmail := EmailRequest{
		From:    "orders@yourdomain.com",
		To:      "customer@example.com",
		Subject: "Order Confirmation #12345",
		HTMLBody: `
			<h1>Order Confirmed!</h1>
			<p>Your order #12345 has been confirmed.</p>
			<h2>Order Details:</h2>
			<ul>
				<li>Product: Premium Widget</li>
				<li>Quantity: 2</li>
				<li>Total: $49.99</li>
			</ul>
			<p><a href="https://example.com/orders/12345">View Order</a></p>
		`,
		TextBody: "Order Confirmed! Your order #12345 has been confirmed. View: https://example.com/orders/12345",
	}

	resp, err = client.SendEmail(transactionalEmail)
	if err != nil {
		fmt.Printf("Error sending transactional email: %v\n", err)
	} else {
		fmt.Printf("Success! Email ID: %s, Remaining quota: %d\n", resp.EmailID, resp.Remaining)
	}
}
