# UJUMBE API Documentation

## Overview
UJUMBE is a powerful email delivery platform built for the Kenyan market with Go, PostgreSQL, Kafka, AWS SES, and Svelte.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All protected endpoints require an API key in the header:
```
X-API-Key: your-api-key-here
```
Or:
```
Authorization: Bearer your-api-key-here
```

## Endpoints

### Authentication

#### Register
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "api_key": "uuid-api-key",
  "user_id": "user-uuid"
}
```

#### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "success": true,
  "token": "api-key",
  "api_key": "api-key",
  "message": "Login successful"
}
```

### Email Operations

#### Send Email
```http
POST /api/v1/emails/send
```

**Request Body (Direct):**
```json
{
  "from": "sender@yourdomain.com",
  "to": "recipient@example.com",
  "subject": "Hello!",
  "html_body": "<h1>Hello World</h1>",
  "text_body": "Hello World"
}
```

**Request Body (Template):**
```json
{
  "from": "sender@yourdomain.com",
  "to": "recipient@example.com",
  "template_id": "uuid-template-id",
  "template_data": {
    "name": "John",
    "url": "https://example.com"
  }
}
```

**Response:**
```json
{
  "success": true,
  "email_id": "email-uuid",
  "message": "Email queued successfully",
  "remaining": 99
}
```

#### Get Email Status
```http
GET /api/v1/emails/{email_id}
```

**Response:**
```json
{
  "id": "email-uuid",
  "user_id": "user-uuid",
  "from_email": "sender@yourdomain.com",
  "to_email": "recipient@example.com",
  "subject": "Hello!",
  "status": "sent",
  "ses_message_id": "ses-message-id",
  "sent_at": "2024-01-01T12:00:00Z",
  "created_at": "2024-01-01T12:00:00Z"
}
```

#### Get Email Logs
```http
GET /api/v1/emails/{email_id}/logs
```

**Response:**
```json
[
  {
    "id": "log-uuid",
    "email_id": "email-uuid",
    "event_type": "created",
    "created_at": "2024-01-01T12:00:00Z"
  },
  {
    "id": "log-uuid",
    "email_id": "email-uuid",
    "event_type": "sent",
    "created_at": "2024-01-01T12:00:05Z"
  }
]
```

### Templates

#### Create Template
```http
POST /api/v1/templates
```

**Request Body:**
```json
{
  "name": "Welcome Email",
  "subject": "Welcome {{name}}!",
  "html_body": "<h1>Welcome {{name}}!</h1><p>Click <a href='{{url}}'>here</a></p>",
  "text_body": "Welcome {{name}}! Visit: {{url}}"
}
```

**Response:**
```json
{
  "id": "template-uuid",
  "user_id": "user-uuid",
  "name": "Welcome Email",
  "subject": "Welcome {{name}}!",
  "html_body": "...",
  "text_body": "...",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### List Templates
```http
GET /api/v1/templates
```

#### Get Template
```http
GET /api/v1/templates/{template_id}
```

#### Delete Template
```http
DELETE /api/v1/templates/{template_id}
```

### Domains

#### Add Domain
```http
POST /api/v1/domains
```

**Request Body:**
```json
{
  "domain": "yourdomain.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Domain added. Please verify ownership by adding DNS records.",
  "domain_id": "domain-uuid",
  "verification_token": "token"
}
```

#### List Domains
```http
GET /api/v1/domains
```

#### Verify Domain
```http
POST /api/v1/domains/{domain_id}/verify
```

#### Delete Domain
```http
DELETE /api/v1/domains/{domain_id}
```

### Analytics & Quota

#### Get Analytics
```http
GET /api/v1/analytics
```

**Response:**
```json
{
  "total_emails_sent": 50,
  "total_emails_failed": 2,
  "total_emails_pending": 3,
  "free_emails_remaining": 45,
  "paid_emails_balance": 100,
  "success_rate": 96.15
}
```

#### Get Quota
```http
GET /api/v1/quota
```

**Response:**
```json
{
  "id": "quota-uuid",
  "user_id": "user-uuid",
  "free_emails_remaining": 95,
  "paid_emails_balance": 0,
  "monthly_reset_date": "2024-02-01",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

## Error Responses

All endpoints return errors in this format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid API key
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Rate Limiting

- Free tier: 100 emails per month
- Each email deducted from free quota first, then paid balance
- Free quota resets on the first day of each month

## Quota System

1. **Free Emails**: Every user gets 100 free emails monthly
2. **Paid Emails**: Top up via M-Pesa for pay-as-you-go pricing
3. **Email Deduction**: Free emails are used first, then paid balance

## Template Variables

Use `{{variable}}` syntax in templates:
```html
<h1>Hello {{name}}!</h1>
<p>Your order #{{order_id}} has been confirmed.</p>
```

Pass values when sending:
```json
{
  "template_data": {
    "name": "John Doe",
    "order_id": "12345"
  }
}
```
