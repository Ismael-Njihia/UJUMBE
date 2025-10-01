# UJUMBE Quick Start Tutorial

This tutorial will guide you through using UJUMBE to send your first email.

## Step 1: Start UJUMBE

Using Docker Compose (easiest):
```bash
make run
```

Wait for all services to start. Check with:
```bash
docker ps
```

## Step 2: Register Your Account

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "you@example.com",
    "password": "your_secure_password"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "api_key": "your-api-key-here",
  "user_id": "user-id-here"
}
```

💡 **Save your API key** - you'll need it for all subsequent requests!

## Step 3: Verify a Domain

Before sending emails, you need to verify a domain:

```bash
curl -X POST http://localhost:8080/api/v1/domains \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{
    "domain": "yourdomain.com"
  }'
```

**Then:**
1. Go to AWS SES Console
2. Find your domain verification records
3. Add the TXT records to your DNS
4. Wait 5-10 minutes for DNS propagation
5. Verify the domain:

```bash
curl -X POST http://localhost:8080/api/v1/domains/{domain-id}/verify \
  -H "X-API-Key: your-api-key-here"
```

## Step 4: Send Your First Email

### Option A: Direct Email

```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{
    "from": "hello@yourdomain.com",
    "to": "recipient@example.com",
    "subject": "Hello from UJUMBE!",
    "html_body": "<h1>Hello!</h1><p>This is my first email from UJUMBE.</p>",
    "text_body": "Hello! This is my first email from UJUMBE."
  }'
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

### Option B: Using a Template

**1. Create a template:**

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{
    "name": "Welcome Email",
    "subject": "Welcome to {{company}}, {{name}}!",
    "html_body": "<h1>Welcome {{name}}!</h1><p>Thanks for joining {{company}}.</p><p><a href=\"{{activation_link}}\">Activate your account</a></p>",
    "text_body": "Welcome {{name}}! Thanks for joining {{company}}. Activate: {{activation_link}}"
  }'
```

**2. Send using the template:**

```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{
    "from": "hello@yourdomain.com",
    "to": "newuser@example.com",
    "template_id": "template-uuid-from-step-1",
    "template_data": {
      "name": "John Doe",
      "company": "ACME Inc",
      "activation_link": "https://example.com/activate/abc123"
    }
  }'
```

## Step 5: Track Email Status

Check if your email was sent:

```bash
curl http://localhost:8080/api/v1/emails/{email-id} \
  -H "X-API-Key: your-api-key-here"
```

**Response:**
```json
{
  "id": "email-uuid",
  "status": "sent",
  "from_email": "hello@yourdomain.com",
  "to_email": "recipient@example.com",
  "subject": "Hello from UJUMBE!",
  "ses_message_id": "ses-message-id",
  "sent_at": "2024-01-01T12:00:05Z",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**Email Status Values:**
- `pending` - Queued for sending
- `sent` - Successfully sent
- `failed` - Failed to send (check `error_message`)

## Step 6: View Email Logs

Get detailed logs for an email:

```bash
curl http://localhost:8080/api/v1/emails/{email-id}/logs \
  -H "X-API-Key: your-api-key-here"
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
    "event_type": "queued",
    "created_at": "2024-01-01T12:00:01Z"
  },
  {
    "id": "log-uuid",
    "email_id": "email-uuid",
    "event_type": "sent",
    "created_at": "2024-01-01T12:00:05Z"
  }
]
```

## Step 7: Check Your Analytics

View your sending statistics:

```bash
curl http://localhost:8080/api/v1/analytics \
  -H "X-API-Key: your-api-key-here"
```

**Response:**
```json
{
  "total_emails_sent": 1,
  "total_emails_failed": 0,
  "total_emails_pending": 0,
  "free_emails_remaining": 99,
  "paid_emails_balance": 0,
  "success_rate": 100.0
}
```

## Step 8: Check Your Quota

See your remaining email quota:

```bash
curl http://localhost:8080/api/v1/quota \
  -H "X-API-Key: your-api-key-here"
```

**Response:**
```json
{
  "free_emails_remaining": 99,
  "paid_emails_balance": 0,
  "monthly_reset_date": "2024-02-01"
}
```

## Using the Dashboard

UJUMBE also includes a web dashboard for easier management.

**Access it at:** http://localhost:3000

**To start the frontend:**
```bash
cd frontend
npm install
npm run dev
```

The dashboard lets you:
- 📧 Send emails with a form
- 📝 Create and manage templates
- 🌐 Add and verify domains
- 📊 View analytics
- 💰 Check quota

## Common Tasks

### List All Templates
```bash
curl http://localhost:8080/api/v1/templates \
  -H "X-API-Key: your-api-key-here"
```

### List All Domains
```bash
curl http://localhost:8080/api/v1/domains \
  -H "X-API-Key: your-api-key-here"
```

### Delete a Template
```bash
curl -X DELETE http://localhost:8080/api/v1/templates/{template-id} \
  -H "X-API-Key: your-api-key-here"
```

## Quota Management

**Free Tier:**
- 100 emails per month
- Resets on the 1st of each month
- Perfect for testing and small projects

**Pay-as-you-go:**
- Top up with M-Pesa (coming soon)
- KES 1 per email (configurable)
- No monthly commitment

## Best Practices

1. **Always verify your domains** before sending production emails
2. **Use templates** for consistent branding and easier management
3. **Monitor analytics** to track delivery success rates
4. **Check logs** if emails fail to send
5. **Test with small batches** before sending to many recipients
6. **Keep your API key secure** - never commit it to version control

## Troubleshooting

### Email shows as "pending" for too long
- Check if the worker is running: `docker ps | grep worker`
- View worker logs: `docker logs ujumbe-worker`

### Email failed to send
- Check email logs for error details
- Verify the sender domain is verified in AWS SES
- Ensure AWS SES is not in sandbox mode (or recipient is verified)

### "Invalid API key" error
- Double-check your API key is correct
- Ensure you're including it in the `X-API-Key` header

### Domain verification failing
- Ensure DNS records are correctly added
- Wait 10-15 minutes for DNS propagation
- Use `dig` or `nslookup` to verify DNS records

## Next Steps

- Read the full [API Documentation](docs/API.md)
- Review the [Setup Guide](docs/SETUP.md) for production deployment
- Explore template variables for dynamic content
- Set up M-Pesa for automatic top-ups
- Integrate UJUMBE with your application

## Need Help?

- 📚 Documentation: Check `docs/` folder
- 🐛 Issues: https://github.com/Ismael-Njihia/UJUMBE/issues
- 📧 Email: support@ujumbe.co.ke

Happy sending! 📧
