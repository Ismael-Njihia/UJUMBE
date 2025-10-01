# UJUMBE - Project Summary

## Overview
UJUMBE is a comprehensive email delivery platform specifically designed for the Kenyan market. It combines modern technology with local payment integration to provide a developer-friendly, reliable, and cost-effective email service.

## Architecture Overview

### Technology Stack
- **Backend**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Message Queue**: Apache Kafka
- **Email Service**: AWS SES
- **Frontend**: Svelte + Vite
- **Containerization**: Docker + Docker Compose
- **Payment**: M-Pesa API

### System Components

#### 1. Backend API (Go)
Located in `backend/cmd/api/`

**Core Services:**
- `auth.go` - User authentication, JWT tokens, API key management
- `email.go` - Email template management, sending, logging
- `domain.go` - Domain verification system
- `billing.go` - M-Pesa integration, transaction management
- `kafka.go` - Message queue integration
- `ses.go` - AWS SES email delivery

**Features:**
- RESTful API with JWT and API key authentication
- Template-based email system with variable substitution
- Domain verification for sender authentication
- Real-time email logging and status tracking
- Analytics and metrics collection
- M-Pesa STK Push for payments
- Quota management (100 free emails/month)

#### 2. Database Layer
Located in `backend/internal/database/`

**Tables:**
- `users` - User accounts, quota, balance
- `api_keys` - API keys for authentication
- `email_templates` - Reusable email templates
- `domains` - Verified sender domains
- `email_logs` - Email delivery logs
- `transactions` - Billing transactions
- `analytics` - Daily metrics

**Features:**
- Automatic schema creation
- Indexed for performance
- Parameterized queries (SQL injection prevention)
- Support for JSONB (template variables)

#### 3. Message Queue (Kafka)
Located in `backend/internal/kafka/`

**Purpose:**
- Asynchronous email processing
- Decouples API from email delivery
- Ensures reliability and scalability
- Handles high volumes efficiently

**Flow:**
1. API receives email request
2. Validates and logs to database
3. Publishes to Kafka queue
4. Consumer picks up and sends via SES
5. Updates log with delivery status

#### 4. Frontend Dashboard (Svelte)
Located in `frontend/src/`

**Pages:**
- `Login.svelte` - Authentication
- `Dashboard.svelte` - Overview with stats
- `Templates.svelte` - Email template management
- `Domains.svelte` - Domain verification
- `APIKeys.svelte` - API key management
- `EmailLogs.svelte` - Real-time logs
- `Analytics.svelte` - Charts and metrics
- `Billing.svelte` - M-Pesa top-up

**Features:**
- Modern, responsive UI
- Real-time data updates
- Interactive charts
- Easy-to-use forms
- Mobile-friendly

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token

#### User Management
- `GET /api/v1/user` - Get current user info
- `GET /api/v1/user/stats` - Dashboard statistics

#### API Keys
- `GET /api/v1/api-keys` - List API keys
- `POST /api/v1/api-keys` - Create new API key
- `DELETE /api/v1/api-keys/{id}` - Revoke API key

#### Email Templates
- `GET /api/v1/templates` - List templates
- `POST /api/v1/templates` - Create template
- `GET /api/v1/templates/{id}` - Get template
- `PUT /api/v1/templates/{id}` - Update template
- `DELETE /api/v1/templates/{id}` - Delete template

#### Email Sending
- `POST /api/v1/emails/send` - Send email (with or without template)
- `GET /api/v1/emails/logs` - Get email logs

#### Domain Verification
- `GET /api/v1/domains` - List domains
- `POST /api/v1/domains` - Add domain
- `GET /api/v1/domains/{id}` - Get domain details
- `POST /api/v1/domains/{id}/verify` - Verify domain
- `DELETE /api/v1/domains/{id}` - Delete domain

#### Analytics
- `GET /api/v1/analytics` - Get analytics data

#### Billing
- `GET /api/v1/billing/transactions` - List transactions
- `POST /api/v1/billing/topup` - Initiate M-Pesa top-up
- `POST /api/v1/mpesa/callback` - M-Pesa callback webhook

## Key Features

### 1. Developer-Friendly API
- Clear, RESTful endpoints
- JWT and API key authentication
- Comprehensive error messages
- Well-documented
- Client libraries examples (Python, Node.js, PHP)

### 2. Email Templates
- Create reusable templates
- Variable substitution ({{variable}})
- HTML and plain text versions
- Preview support
- Easy management interface

### 3. Domain Verification
- Add custom sender domains
- DNS-based verification
- Improves deliverability
- DKIM support

### 4. Real-time Logging
- Track every email
- Status updates (queued, sent, failed, bounced)
- Error messages
- Message IDs
- Searchable logs

### 5. Analytics
- Daily email metrics
- Success/failure rates
- Bounce tracking
- Visual charts
- Exportable data

### 6. Quota Management
- 100 free emails per month
- Automatic quota reset
- Pay-as-you-go pricing (KES 1.00/email)
- Balance tracking
- Usage notifications

### 7. M-Pesa Integration
- STK Push for easy payments
- Real-time callback processing
- Transaction history
- Automatic credit addition
- Secure payment flow

### 8. Security
- Bcrypt password hashing
- JWT token authentication
- API key authentication
- Environment-based secrets
- SQL injection prevention
- CORS protection
- Domain verification

### 9. Scalability
- Kafka for async processing
- Horizontal scaling support
- Database indexing
- Connection pooling
- Stateless API design

### 10. Reliability
- Error handling
- Retry mechanisms
- Transaction tracking
- Monitoring support
- Health checks

## Deployment Options

### Docker Compose (Recommended)
```bash
docker-compose up -d
```
Includes:
- PostgreSQL database
- Kafka + Zookeeper
- API server
- All networking configured

### Manual Deployment
- Systemd service
- Nginx reverse proxy
- Let's Encrypt SSL
- Separate components

### Cloud Deployment
- AWS/GCP/Azure compatible
- Container orchestration (Kubernetes)
- Managed services (RDS, MSK)
- Load balancing

## Configuration

### Environment Variables
All configuration via `.env` file:
- Database connection
- Kafka brokers
- AWS credentials
- M-Pesa credentials
- JWT secret
- CORS origins

### Customization
- Email templates
- Pricing (configurable per email cost)
- Free quota amount
- Domain verification requirements

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL (if not using Docker)
- Kafka (if not using Docker)

### Quick Start
```bash
# Clone repository
git clone https://github.com/Ismael-Njihia/UJUMBE.git
cd UJUMBE

# Configure environment
cp .env.example .env
# Edit .env with your credentials

# Start with Docker
docker-compose up -d

# Or run manually
make run          # Backend
make frontend-dev # Frontend
```

### Testing
```bash
make test      # Run tests
make lint      # Run linters
make build     # Build binary
```

## Documentation

### Available Docs
- `README.md` - Main documentation
- `docs/README.md` - Comprehensive guide
- `docs/API_EXAMPLES.md` - API usage examples
- `docs/DEPLOYMENT.md` - Deployment guide
- `CONTRIBUTING.md` - Contributing guidelines
- `CHANGELOG.md` - Version history

### Code Structure
```
UJUMBE/
├── backend/
│   ├── cmd/api/          # Main application
│   ├── internal/         # Internal packages
│   │   ├── auth/         # Authentication
│   │   ├── email/        # Email service
│   │   ├── domain/       # Domain verification
│   │   ├── billing/      # Billing & payments
│   │   ├── kafka/        # Message queue
│   │   ├── database/     # Database layer
│   │   └── models/       # Data models
│   └── pkg/              # Shared packages
│       ├── config/       # Configuration
│       ├── ses/          # AWS SES client
│       └── mpesa/        # M-Pesa client
├── frontend/
│   └── src/
│       ├── routes/       # Svelte pages
│       ├── lib/          # Shared code
│       └── components/   # Reusable components
├── database/
│   └── migrations/       # Database migrations
├── docs/                 # Documentation
├── .github/workflows/    # CI/CD
└── docker-compose.yml    # Docker setup
```

## Future Enhancements

### Planned Features
- Email webhooks for delivery notifications
- Email scheduling
- A/B testing for templates
- Advanced analytics (open rates, click tracking)
- Multi-language support
- Mobile app
- Additional payment methods
- Email list management
- Automated email campaigns
- Rate limiting middleware
- Redis caching
- Email bounce handling
- Unsubscribe management

### Potential Integrations
- Slack notifications
- Webhook callbacks
- SMS notifications (via M-Pesa or other services)
- CDN for attachments
- Email validation API
- DMARC/DKIM automation

## Performance Metrics

### Expected Performance
- API Response Time: < 100ms
- Email Queue Time: < 1 second
- Email Delivery: 2-5 seconds (via SES)
- Dashboard Load: < 2 seconds
- Database Queries: < 50ms

### Scalability
- Handles 10,000+ emails/hour
- Horizontal scaling support
- Multiple consumer workers
- Database read replicas
- CDN for static assets

## Support & Community

### Getting Help
- GitHub Issues
- Documentation
- API Examples
- Community Forum (planned)

### Contributing
See `CONTRIBUTING.md` for:
- Code standards
- Pull request process
- Development setup
- Testing guidelines

## License
MIT License - See `LICENSE` file

## Credits
Built for the Kenyan market with ❤️

---

**Version**: 1.0.0  
**Last Updated**: January 2024  
**Status**: Production Ready
