# API Usage Examples

This document provides practical examples of using the UJUMBE email API.

## Authentication

All API requests require authentication using either:
- JWT token (for dashboard/user actions)
- API key (for programmatic email sending)

## Example 1: Complete Workflow

### Step 1: Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "developer@example.com",
    "password": "SecurePass123!",
    "name": "John Developer"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "developer@example.com",
    "password": "SecurePass123!"
  }'

# Response includes JWT token
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { ... }
}
```

### Step 2: Create an API Key

```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Production API Key"
  }'

# Response
{
  "id": "key-uuid",
  "key": "ujumbe_abc123xyz...",
  "name": "Production API Key"
}
```

### Step 3: Create an Email Template

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Welcome Email",
    "subject": "Welcome to {{company_name}}, {{first_name}}!",
    "html_body": "<html><body><h1>Welcome {{first_name}} {{last_name}}!</h1><p>Thank you for joining {{company_name}}.</p></body></html>",
    "text_body": "Welcome {{first_name}} {{last_name}}! Thank you for joining {{company_name}}.",
    "variables": {
      "first_name": "User'\''s first name",
      "last_name": "User'\''s last name",
      "company_name": "Your company name"
    }
  }'
```

### Step 4: Send Email Using Template

```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: ujumbe_abc123xyz..." \
  -d '{
    "template_id": "template-uuid-from-step-3",
    "from": "noreply@yourdomain.com",
    "to": "customer@example.com",
    "variables": {
      "first_name": "Jane",
      "last_name": "Doe",
      "company_name": "ACME Corp"
    }
  }'

# Response
{
  "message": "Email queued successfully",
  "status": "queued"
}
```

## Example 2: Send Simple Email Without Template

```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: ujumbe_abc123xyz..." \
  -d '{
    "from": "noreply@yourdomain.com",
    "to": "customer@example.com",
    "subject": "Your Order Confirmation",
    "html_body": "<h1>Order Confirmed</h1><p>Your order #12345 has been confirmed.</p>",
    "text_body": "Order Confirmed. Your order #12345 has been confirmed."
  }'
```

## Example 3: Add and Verify Domain

```bash
# Add domain
curl -X POST http://localhost:8080/api/v1/domains \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "domain": "yourdomain.com"
  }'

# Response includes verification code
{
  "id": "domain-uuid",
  "domain": "yourdomain.com",
  "verification_code": "ujumbe-verify-abc123",
  "is_verified": false
}

# Add TXT record to your DNS:
# Type: TXT
# Name: _ujumbe-verify
# Value: ujumbe-verify-abc123

# After adding DNS record, verify:
curl -X POST http://localhost:8080/api/v1/domains/domain-uuid/verify \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Example 4: Check Email Logs

```bash
# Get recent email logs
curl -X GET "http://localhost:8080/api/v1/emails/logs?limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response
[
  {
    "id": "log-uuid",
    "from_email": "noreply@yourdomain.com",
    "to_email": "customer@example.com",
    "subject": "Welcome to ACME Corp, Jane!",
    "status": "sent",
    "created_at": "2024-01-15T10:30:00Z",
    "sent_at": "2024-01-15T10:30:05Z"
  }
]
```

## Example 5: Get Analytics

```bash
# Get analytics for last 30 days
curl -X GET "http://localhost:8080/api/v1/analytics?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response
[
  {
    "date": "2024-01-15",
    "emails_sent": 150,
    "failed": 2,
    "bounced": 1
  },
  {
    "date": "2024-01-14",
    "emails_sent": 120,
    "failed": 0,
    "bounced": 0
  }
]
```

## Example 6: Top Up with M-Pesa

```bash
# Initiate M-Pesa payment
curl -X POST http://localhost:8080/api/v1/billing/topup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 500,
    "phone_number": "254712345678"
  }'

# Response
{
  "transaction_id": "txn-uuid",
  "mpesa_response": {
    "CheckoutRequestID": "ws_CO_123456",
    "ResponseDescription": "Success. Request accepted for processing"
  },
  "message": "Payment initiated. Please complete on your phone"
}

# User will receive STK push on their phone
# After payment completion, webhook updates the transaction
```

## Example 7: Check Dashboard Stats

```bash
curl -X GET http://localhost:8080/api/v1/user/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response
{
  "email_quota": 100,
  "emails_sent": 45,
  "balance": 150.00,
  "total_sent": 450,
  "total_failed": 5,
  "total_bounced": 2,
  "quota_remaining": 55
}
```

## Integration Examples

### Python Example

```python
import requests

API_KEY = "ujumbe_abc123xyz..."
BASE_URL = "http://localhost:8080/api/v1"

def send_email(to_email, subject, html_body):
    headers = {
        "X-API-Key": API_KEY,
        "Content-Type": "application/json"
    }
    
    data = {
        "from": "noreply@yourdomain.com",
        "to": to_email,
        "subject": subject,
        "html_body": html_body,
        "text_body": html_body  # Simple conversion
    }
    
    response = requests.post(
        f"{BASE_URL}/emails/send",
        headers=headers,
        json=data
    )
    
    return response.json()

# Usage
result = send_email(
    "customer@example.com",
    "Welcome!",
    "<h1>Welcome to our platform!</h1>"
)
print(result)
```

### Node.js Example

```javascript
const axios = require('axios');

const API_KEY = 'ujumbe_abc123xyz...';
const BASE_URL = 'http://localhost:8080/api/v1';

async function sendEmail(toEmail, subject, htmlBody) {
  try {
    const response = await axios.post(
      `${BASE_URL}/emails/send`,
      {
        from: 'noreply@yourdomain.com',
        to: toEmail,
        subject: subject,
        html_body: htmlBody,
        text_body: htmlBody
      },
      {
        headers: {
          'X-API-Key': API_KEY,
          'Content-Type': 'application/json'
        }
      }
    );
    
    return response.data;
  } catch (error) {
    console.error('Error sending email:', error.response.data);
    throw error;
  }
}

// Usage
sendEmail(
  'customer@example.com',
  'Welcome!',
  '<h1>Welcome to our platform!</h1>'
).then(result => {
  console.log('Email sent:', result);
});
```

### PHP Example

```php
<?php

function sendEmail($toEmail, $subject, $htmlBody) {
    $apiKey = 'ujumbe_abc123xyz...';
    $baseUrl = 'http://localhost:8080/api/v1';
    
    $data = [
        'from' => 'noreply@yourdomain.com',
        'to' => $toEmail,
        'subject' => $subject,
        'html_body' => $htmlBody,
        'text_body' => $htmlBody
    ];
    
    $ch = curl_init($baseUrl . '/emails/send');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'X-API-Key: ' . $apiKey,
        'Content-Type: application/json'
    ]);
    
    $response = curl_exec($ch);
    curl_close($ch);
    
    return json_decode($response, true);
}

// Usage
$result = sendEmail(
    'customer@example.com',
    'Welcome!',
    '<h1>Welcome to our platform!</h1>'
);

print_r($result);
?>
```

## Error Handling

All API responses follow this structure:

**Success Response:**
```json
{
  "message": "Success message",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- 200: Success
- 201: Created
- 202: Accepted (for queued emails)
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden (quota exceeded)
- 404: Not Found
- 500: Internal Server Error

## Rate Limiting

- Free tier: 100 emails per month
- After free tier: Pay-as-you-go (KES 1.00 per email)
- No hard rate limits on API requests
- Quota resets on the 1st of each month
