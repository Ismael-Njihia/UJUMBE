# UJUMBE - Email Delivery Platform for Kenya

A secure, developer-first email platform built for Kenya using Go, PostgreSQL, Kafka, and AWS SES with a Svelte dashboard. It offers APIs with template IDs, verified sender domains, real-time logs, analytics, and M-Pesa billing. Users get 100 free emails monthly then pay-as-you-go, ensuring fast, reliable delivery.

## Features

### Core Features
- ✅ **Developer-Friendly API** - RESTful API with JWT and API key authentication
- ✅ **Email Templates** - Create and manage reusable email templates with variables
- ✅ **Domain Verification** - Verify custom sender domains for improved deliverability
- ✅ **Real-time Logs** - Track every email with detailed status information
- ✅ **Analytics Dashboard** - Monitor email performance with charts and metrics
- ✅ **Quota Management** - 100 free emails per month, auto-reset

### Payment & Billing
- 💳 **M-Pesa Integration** - Easy top-up with M-Pesa STK Push
- 💰 **Pay-as-you-go** - KES 1.00 per email after free tier
- 📊 **Transaction History** - Track all payments and credits

### Technical Stack
- **Backend**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Message Queue**: Apache Kafka
- **Email Service**: AWS SES
- **Frontend**: Svelte + Vite
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- AWS SES credentials (for sending emails)
- M-Pesa API credentials (for payments)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/Ismael-Njihia/UJUMBE.git
cd UJUMBE
```

2. **Configure environment variables**
```bash
cp .env.example .env
# Edit .env with your credentials
```

3. **Start the services**
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Local Development

#### Backend Development

```bash
# Install Go dependencies
cd backend
go mod download

# Run the API server
go run cmd/api/*.go
```

#### Frontend Development

```bash
# Install npm dependencies
cd frontend
npm install

# Start development server
npm run dev
```

The dashboard will be available at `http://localhost:5173`

## API Documentation

### Authentication

#### Register
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Doe"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password"
}

Response:
{
  "token": "jwt_token_here",
  "user": { ... }
}
```

### Email Templates

#### Create Template
```bash
POST /api/v1/templates
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "name": "Welcome Email",
  "subject": "Welcome {{name}}!",
  "html_body": "<h1>Hello {{name}}</h1><p>Welcome to our platform!</p>",
  "text_body": "Hello {{name}}, Welcome to our platform!",
  "variables": {
    "name": "User's name"
  }
}
```

#### Get Templates
```bash
GET /api/v1/templates
Authorization: Bearer {jwt_token}
```

### Sending Emails

#### Send Email with Template
```bash
POST /api/v1/emails/send
X-API-Key: {your_api_key}
Content-Type: application/json

{
  "template_id": "template_uuid",
  "from": "noreply@yourdomain.com",
  "to": "recipient@example.com",
  "variables": {
    "name": "John Doe"
  }
}
```

#### Send Email without Template
```bash
POST /api/v1/emails/send
X-API-Key: {your_api_key}
Content-Type: application/json

{
  "from": "noreply@yourdomain.com",
  "to": "recipient@example.com",
  "subject": "Test Email",
  "html_body": "<h1>Hello</h1><p>This is a test email.</p>",
  "text_body": "Hello, This is a test email."
}
```

### Domain Verification

#### Add Domain
```bash
POST /api/v1/domains
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "domain": "yourdomain.com"
}

Response:
{
  "id": "domain_uuid",
  "domain": "yourdomain.com",
  "verification_code": "verification_code_here",
  "is_verified": false
}
```

Add the verification code as a TXT record to your domain's DNS settings.

#### Verify Domain
```bash
POST /api/v1/domains/{domain_id}/verify
Authorization: Bearer {jwt_token}
```

### API Keys

#### Create API Key
```bash
POST /api/v1/api-keys
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "name": "Production Server"
}

Response:
{
  "id": "key_uuid",
  "key": "ujumbe_xxxxxxxxxxxxx",
  "name": "Production Server"
}
```

**Important**: Store the API key securely. It won't be shown again.

### Email Logs

#### Get Email Logs
```bash
GET /api/v1/emails/logs?limit=50&offset=0
Authorization: Bearer {jwt_token}
```

### Analytics

#### Get Analytics
```bash
GET /api/v1/analytics?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer {jwt_token}
```

### Billing

#### Top Up with M-Pesa
```bash
POST /api/v1/billing/topup
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "amount": 100,
  "phone_number": "254712345678"
}
```

You'll receive an STK push on your phone to complete the payment.

#### Get Transactions
```bash
GET /api/v1/billing/transactions?limit=50&offset=0
Authorization: Bearer {jwt_token}
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `KAFKA_BROKERS` | Kafka broker addresses | localhost:9092 |
| `AWS_REGION` | AWS region for SES | us-east-1 |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - |
| `MPESA_CONSUMER_KEY` | M-Pesa consumer key | - |
| `MPESA_CONSUMER_SECRET` | M-Pesa consumer secret | - |
| `MPESA_SHORTCODE` | M-Pesa business shortcode | - |
| `MPESA_PASSKEY` | M-Pesa passkey | - |
| `JWT_SECRET` | Secret for JWT tokens | - |
| `APP_PORT` | API server port | 8080 |
| `FREE_EMAIL_QUOTA` | Free emails per month | 100 |

## Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Svelte    │─────▶│   Go API     │─────▶│ PostgreSQL  │
│  Dashboard  │      │   (Backend)  │      │  (Database) │
└─────────────┘      └──────────────┘      └─────────────┘
                            │
                            │
                     ┌──────┴──────┐
                     │             │
                     ▼             ▼
              ┌───────────┐  ┌──────────┐
              │   Kafka   │  │ AWS SES  │
              │  (Queue)  │  │ (Email)  │
              └───────────┘  └──────────┘
                     │
                     ▼
              ┌───────────┐
              │  Consumer │
              │  Workers  │
              └───────────┘
```

### Email Flow

1. User sends email request to API
2. API validates quota and creates log entry
3. Email job is queued in Kafka
4. Consumer worker picks up the job
5. Email is sent via AWS SES
6. Log is updated with status
7. Analytics are recorded

## Security

- ✅ JWT-based authentication for dashboard
- ✅ API key authentication for programmatic access
- ✅ Password hashing with bcrypt
- ✅ Domain verification for sender authentication
- ✅ CORS protection
- ✅ SQL injection prevention with parameterized queries
- ✅ Rate limiting (via quota system)

## Monitoring & Logs

- Real-time email status tracking
- Daily analytics aggregation
- Detailed error logging
- Transaction history

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/Ismael-Njihia/UJUMBE/issues
- Email: support@ujumbe.ke (placeholder)

## License

MIT License - feel free to use for commercial or personal projects.

## Roadmap

- [ ] Email webhooks for delivery notifications
- [ ] Email scheduling
- [ ] A/B testing for templates
- [ ] Advanced analytics (open rates, click tracking)
- [ ] Multi-language support
- [ ] Mobile app
- [ ] Additional payment methods
- [ ] Email list management
- [ ] Automated email campaigns

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

---

Built with ❤️ for the Kenyan market
