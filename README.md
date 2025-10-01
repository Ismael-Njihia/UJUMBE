# UJUMBE 📧

A secure, developer-first email delivery platform built for the Kenyan market. UJUMBE provides powerful APIs with template support, custom verified sender domains, real-time logs, analytics, and M-Pesa billing integration.

## Features

✨ **Core Features**
- 🚀 Fast and reliable email delivery via AWS SES
- 📝 Email templates with variable substitution
- 🔐 Custom verified sender domains
- 📊 Real-time email logs and tracking
- 📈 Comprehensive analytics dashboard
- 💳 M-Pesa integration for pay-as-you-go billing
- 🆓 100 free emails monthly per user

🛠 **Technical Stack**
- **Backend**: Go 1.21
- **Database**: PostgreSQL 15
- **Message Queue**: Apache Kafka
- **Email Service**: AWS SES
- **Frontend**: Svelte + Vite
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites
- Docker and Docker Compose
- AWS account with SES configured
- (Optional) M-Pesa API credentials

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ismael-Njihia/UJUMBE.git
   cd UJUMBE
   ```

2. **Configure environment variables**
   ```bash
   cp backend/.env.example backend/.env
   # Edit backend/.env with your AWS credentials
   ```

3. **Start the services**
   ```bash
   make run
   ```

   This will start:
   - PostgreSQL database on port 5432
   - Kafka on port 9092
   - Backend API on port 8080
   - Email worker service

4. **Access the services**
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health
   - API Documentation: [docs/API.md](docs/API.md)

### Development Setup

**Backend Development**
```bash
cd backend
go mod download
go run cmd/server/main.go
```

**Worker Development**
```bash
cd backend
go run cmd/worker/main.go
```

**Frontend Development**
```bash
cd frontend
npm install
npm run dev
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────┐
│   Svelte    │────▶│   Go API     │────▶│  Kafka  │
│  Dashboard  │     │   Server     │     │  Queue  │
└─────────────┘     └──────────────┘     └─────────┘
                           │                    │
                           ▼                    ▼
                    ┌──────────────┐     ┌──────────┐
                    │  PostgreSQL  │     │  Worker  │
                    │   Database   │     │ Service  │
                    └──────────────┘     └──────────┘
                                               │
                                               ▼
                                        ┌──────────┐
                                        │ AWS SES  │
                                        └──────────┘
```

## API Usage

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123"}'
```

### Send an email
```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "from": "sender@yourdomain.com",
    "to": "recipient@example.com",
    "subject": "Hello from UJUMBE",
    "html_body": "<h1>Hello!</h1><p>This is a test email.</p>"
  }'
```

### Create a template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "Welcome Email",
    "subject": "Welcome {{name}}!",
    "html_body": "<h1>Welcome {{name}}!</h1>"
  }'
```

For complete API documentation, see [docs/API.md](docs/API.md).

## Database Schema

The system uses PostgreSQL with the following main tables:
- `users` - User accounts and API keys
- `user_quotas` - Email quotas (free and paid)
- `verified_domains` - Custom sender domains
- `email_templates` - Reusable email templates
- `emails` - Email tracking and status
- `email_logs` - Real-time delivery logs
- `transactions` - M-Pesa payment records

See [database/migrations/001_init_schema.sql](database/migrations/001_init_schema.sql) for the complete schema.

## Makefile Commands

```bash
make help          # Show all available commands
make build         # Build Docker images
make run           # Start all services
make stop          # Stop all services
make logs          # View logs from all services
make clean         # Remove containers and volumes
make test          # Run tests
make dev           # Run backend in dev mode
make worker-dev    # Run worker in dev mode
```

## Configuration

Key environment variables in `backend/.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ujumbe

# AWS SES
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=email_jobs

# Application
FREE_EMAIL_QUOTA=100
EMAIL_PRICE_PER_UNIT=1.0
```

## Quota System

- **Free Tier**: 100 emails per month
- **Pay-as-you-go**: Top up via M-Pesa
- **Auto-reset**: Free quota resets on the 1st of each month
- **Priority**: Free emails used first, then paid balance

## Security Features

- 🔒 API key authentication
- 🛡️ Password hashing with bcrypt
- ✅ Domain verification via AWS SES
- 🔐 SQL injection prevention with prepared statements
- 🌐 CORS configuration
- 📝 Request validation

## M-Pesa Integration

UJUMBE integrates with Safaricom's M-Pesa for seamless billing:
1. Users request a top-up via the API
2. M-Pesa STK push sent to user's phone
3. User completes payment
4. Credits automatically added to account
5. Real-time transaction tracking

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the test suite:
```bash
cd backend
go test ./...
```

## Production Deployment

For production deployment:
1. Set `ENVIRONMENT=production` in `.env`
2. Use strong JWT secrets and passwords
3. Configure AWS SES production access
4. Set up proper monitoring and logging
5. Use a reverse proxy (nginx) for the API
6. Enable SSL/TLS certificates
7. Configure database backups
8. Set up Kafka cluster for reliability

## Support

- 📧 Email: support@ujumbe.co.ke
- 📚 Documentation: [docs/API.md](docs/API.md)
- 🐛 Issues: [GitHub Issues](https://github.com/Ismael-Njihia/UJUMBE/issues)

## License

MIT License - see LICENSE file for details

## Acknowledgments

Built with ❤️ for the Kenyan developer community

---

**UJUMBE** - *Powerful email delivery for Kenya* 🇰🇪
