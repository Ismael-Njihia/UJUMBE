# UJUMBE - Email Delivery Platform for Kenya

A secure, developer-first email platform built for Kenya using Go, PostgreSQL, Kafka, and AWS SES with a Svelte dashboard. It offers APIs with template IDs, verified sender domains, real-time logs, analytics, and M-Pesa billing. Users get 100 free emails monthly then pay-as-you-go, ensuring fast, reliable delivery.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 🚀 Features

### Core Features
- ✅ **Developer-Friendly API** - RESTful API with JWT and API key authentication
- ✅ **Email Templates** - Create and manage reusable email templates with variables
- ✅ **Domain Verification** - Verify custom sender domains for improved deliverability
- ✅ **Real-time Logs** - Track every email with detailed status information
- ✅ **Analytics Dashboard** - Monitor email performance with charts and metrics
- ✅ **Quota Management** - 100 free emails per month, auto-reset monthly

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

## 📦 Quick Start

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

## 📚 Documentation

For detailed documentation, see [docs/README.md](docs/README.md)

### API Examples

#### Send Email with Template
```bash
POST /api/v1/emails/send
X-API-Key: your_api_key
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

## 🏗️ Architecture

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

## 🔐 Security

- JWT-based authentication for dashboard
- API key authentication for programmatic access
- Password hashing with bcrypt
- Domain verification for sender authentication
- CORS protection
- SQL injection prevention

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 🌟 Support

If you find this project helpful, please give it a ⭐️!

---

Built with ❤️ for the Kenyan market
