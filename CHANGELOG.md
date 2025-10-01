# Changelog

All notable changes to UJUMBE will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

### Added

#### Backend
- Go 1.21+ backend with RESTful API
- PostgreSQL database integration with automatic schema creation
- User authentication with JWT tokens
- API key generation and management for programmatic access
- Email template system with variable substitution
- Domain verification for custom sender domains
- Email queue management using Apache Kafka
- Real-time email logging and status tracking
- Analytics and metrics collection
- M-Pesa STK Push integration for payments
- AWS SES integration for email delivery
- Quota management (100 free emails per month)
- Billing and transaction history
- CORS support for frontend integration

#### Database Schema
- Users table with quota and balance tracking
- API keys table for authentication
- Email templates with JSON variable support
- Domains table with verification codes
- Email logs with detailed status tracking
- Transactions table for billing
- Analytics table for daily metrics
- Comprehensive indexes for performance

#### Frontend
- Modern Svelte dashboard with Vite
- User authentication (login/register)
- Dashboard with statistics overview
- Email template management UI
- Domain verification interface
- API key management
- Real-time email logs viewer
- Analytics charts and metrics
- Billing and M-Pesa top-up interface
- Responsive design

#### Infrastructure
- Docker containerization
- Docker Compose setup with PostgreSQL, Zookeeper, and Kafka
- GitHub Actions CI/CD workflow
- Environment configuration with .env
- Makefile for common tasks
- Quick start script

#### Documentation
- Comprehensive README
- API examples in multiple languages (Python, Node.js, PHP)
- Deployment guide for production
- Contributing guidelines
- MIT License
- Detailed API documentation

### Features

- ✅ Developer-friendly RESTful API
- ✅ JWT and API key authentication
- ✅ Template-based email system
- ✅ Custom domain verification
- ✅ Real-time logging
- ✅ Analytics dashboard
- ✅ M-Pesa payment integration
- ✅ Quota management
- ✅ Asynchronous email processing
- ✅ Secure password hashing
- ✅ CORS protection
- ✅ SQL injection prevention

### Security

- Bcrypt password hashing
- JWT token-based authentication
- API key authentication for programmatic access
- Environment-based secret management
- Parameterized SQL queries
- CORS configuration
- Domain verification for sender authentication

### Performance

- Kafka-based asynchronous email processing
- Database indexing for fast queries
- Connection pooling
- Efficient batch processing
- Horizontal scalability support

## [Unreleased]

### Planned Features
- [ ] Email webhooks for delivery notifications
- [ ] Email scheduling
- [ ] A/B testing for templates
- [ ] Advanced analytics (open rates, click tracking)
- [ ] Multi-language support
- [ ] Mobile app
- [ ] Additional payment methods
- [ ] Email list management
- [ ] Automated email campaigns
- [ ] Rate limiting middleware
- [ ] Redis caching
- [ ] Email bounce handling
- [ ] Unsubscribe management

---

For more details, see the [documentation](docs/README.md).
