# UJUMBE Project Structure

This document describes the structure and organization of the UJUMBE project.

## Directory Structure

```
UJUMBE/
├── backend/                      # Go backend service
│   ├── cmd/                      # Application entry points
│   │   ├── server/              # Main API server
│   │   │   └── main.go
│   │   └── worker/              # Email worker service
│   │       └── main.go
│   ├── internal/                # Internal packages (not importable by other projects)
│   │   ├── api/                 # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── domain_handler.go
│   │   │   ├── email_handler.go
│   │   │   └── template_handler.go
│   │   ├── database/            # Database connection
│   │   │   └── database.go
│   │   ├── kafka/               # Kafka producer/consumer
│   │   │   └── kafka.go
│   │   ├── middleware/          # HTTP middleware
│   │   │   └── auth.go
│   │   ├── models/              # Data models
│   │   │   └── models.go
│   │   ├── services/            # Business logic
│   │   │   ├── auth_service.go
│   │   │   ├── email_service.go
│   │   │   └── services_test.go
│   │   └── ses/                 # AWS SES client
│   │       └── ses.go
│   ├── pkg/                     # Public packages (can be imported)
│   │   └── utils/               # Utility functions
│   │       ├── template.go
│   │       └── template_test.go
│   ├── .env.example             # Environment variables template
│   ├── Dockerfile               # Server container
│   ├── Dockerfile.worker        # Worker container
│   ├── go.mod                   # Go module definition
│   └── go.sum                   # Go dependencies checksums
├── frontend/                    # Svelte frontend
│   ├── src/
│   │   ├── components/          # Reusable components
│   │   ├── lib/                 # Utilities and API client
│   │   │   └── api.js
│   │   ├── routes/              # Page components
│   │   │   ├── Analytics.svelte
│   │   │   ├── Dashboard.svelte
│   │   │   ├── Domains.svelte
│   │   │   ├── Login.svelte
│   │   │   ├── Register.svelte
│   │   │   ├── SendEmail.svelte
│   │   │   └── Templates.svelte
│   │   ├── App.svelte           # Root component
│   │   └── main.js              # Entry point
│   ├── index.html               # HTML template
│   ├── package.json             # Node dependencies
│   ├── svelte.config.js         # Svelte configuration
│   └── vite.config.js           # Vite configuration
├── database/                    # Database files
│   └── migrations/              # SQL migration scripts
│       └── 001_init_schema.sql
├── docs/                        # Documentation
│   ├── API.md                   # API documentation
│   ├── SETUP.md                 # Setup guide
│   └── TUTORIAL.md              # Quick start tutorial
├── examples/                    # Integration examples
│   ├── go-client/               # Go client example
│   │   ├── main.go
│   │   └── README.md
│   └── python-client/           # Python client example
│       ├── ujumbe_client.py
│       ├── requirements.txt
│       └── README.md
├── .gitignore                   # Git ignore rules
├── CONTRIBUTING.md              # Contribution guidelines
├── docker-compose.yml           # Docker services configuration
├── LICENSE                      # MIT License
├── Makefile                     # Build and development commands
└── README.md                    # Main project documentation
```

## Component Description

### Backend (`/backend`)

The backend is written in Go and follows a clean architecture pattern:

#### `cmd/` - Application Entry Points
- **server**: Main API server that handles HTTP requests
- **worker**: Background worker that processes email jobs from Kafka

#### `internal/` - Internal Packages
- **api**: HTTP request handlers (controllers)
- **database**: Database connection and queries
- **kafka**: Message queue producer and consumer
- **middleware**: HTTP middleware (authentication, CORS, etc.)
- **models**: Data structures and request/response types
- **services**: Business logic layer
- **ses**: AWS SES email sending client

#### `pkg/` - Public Packages
- **utils**: Helper functions that could be reused in other projects

### Frontend (`/frontend`)

Built with Svelte and Vite for a modern, reactive user interface:

#### `src/routes/` - Pages
- **Login/Register**: User authentication
- **Dashboard**: Overview of quota and quick actions
- **SendEmail**: Form to send emails
- **Templates**: Manage email templates
- **Domains**: Manage verified domains
- **Analytics**: View sending statistics

#### `src/lib/` - Shared Code
- **api.js**: API client with axios

### Database (`/database`)

SQL migration scripts for PostgreSQL:
- Schema definitions
- Indexes
- Triggers
- Initial data

### Documentation (`/docs`)

Comprehensive documentation:
- **API.md**: Complete API reference
- **SETUP.md**: Installation and configuration
- **TUTORIAL.md**: Quick start guide

### Examples (`/examples`)

Integration examples in multiple languages:
- **Go**: Native Go client
- **Python**: Python client for Django/Flask apps

## Key Technologies

### Backend Stack
- **Go 1.21**: Backend language
- **PostgreSQL 15**: Database
- **Kafka**: Message queue
- **AWS SES**: Email delivery
- **Docker**: Containerization

### Frontend Stack
- **Svelte 4**: UI framework
- **Vite 5**: Build tool
- **Axios**: HTTP client

### Development Tools
- **Docker Compose**: Local development
- **Make**: Build automation
- **Git**: Version control

## Data Flow

```
1. User Request
   ↓
2. API Server (Go)
   ↓
3. Authentication Middleware
   ↓
4. Handler (auth_handler, email_handler, etc.)
   ↓
5. Service Layer (business logic)
   ↓
6. Database (PostgreSQL) OR Kafka Queue
   ↓
7. Worker (consumes from Kafka)
   ↓
8. AWS SES (sends email)
   ↓
9. Update Database (email status)
```

## Key Concepts

### Email Sending Flow
1. User submits email via API
2. Quota is checked and deducted
3. Email record created in database
4. Job sent to Kafka queue
5. Worker picks up job
6. Email sent via AWS SES
7. Status updated in database

### Authentication Flow
1. User registers (POST /auth/register)
2. API key generated and stored
3. User includes API key in subsequent requests
4. Middleware validates API key
5. User ID attached to request context

### Template System
1. User creates template with variables ({{name}})
2. Template stored in database
3. When sending, template loaded
4. Variables replaced with actual data
5. Final email sent

## Environment Configuration

### Backend Environment Variables
- `DB_*`: Database connection
- `AWS_*`: AWS SES credentials
- `KAFKA_*`: Kafka configuration
- `SERVER_PORT`: API server port
- `JWT_SECRET`: Authentication secret

### Frontend Environment
- API URL configured in `src/lib/api.js`
- Can be overridden per environment

## Building and Running

### Development
```bash
# Backend
cd backend && go run cmd/server/main.go

# Worker
cd backend && go run cmd/worker/main.go

# Frontend
cd frontend && npm run dev
```

### Production
```bash
# Using Docker Compose
make build
make run

# Manual build
cd backend && go build -o server cmd/server/main.go
cd backend && go build -o worker cmd/worker/main.go
cd frontend && npm run build
```

## Testing

### Backend Tests
```bash
cd backend
go test ./...                    # Run all tests
go test ./pkg/utils -v          # Test utils package
go test ./internal/services -v  # Test services
```

### Frontend Tests
```bash
cd frontend
npm test                        # Run frontend tests
```

## Adding New Features

### New API Endpoint
1. Add handler to `internal/api/`
2. Add route in `cmd/server/main.go`
3. Implement service logic in `internal/services/`
4. Update models if needed in `internal/models/`
5. Add tests
6. Update API documentation

### New Frontend Page
1. Create component in `src/routes/`
2. Add route in `src/App.svelte`
3. Add API calls in `src/lib/api.js` if needed
4. Update navigation

## Code Organization Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Direction**: Dependencies flow inward (handlers → services → database)
3. **No Circular Dependencies**: Packages don't import each other in circles
4. **Testability**: Services and utilities are easy to test
5. **Documentation**: Every exported function has a comment

## Security Considerations

- API keys stored hashed in database
- Passwords hashed with bcrypt
- SQL injection prevented with parameterized queries
- CORS configured appropriately
- Environment variables for secrets

## Performance Optimization

- Database indexes on frequently queried columns
- Connection pooling for database
- Kafka for async processing
- Batch operations where possible
- Caching strategies (to be implemented)

## Monitoring and Logging

- Structured logging throughout
- Email status tracking in database
- Error logging with context
- Metrics collection (to be implemented)

## Future Enhancements

See CONTRIBUTING.md for areas needing contribution:
- M-Pesa integration
- More comprehensive testing
- Webhook notifications
- Email scheduling
- Attachment support

## Questions?

- Check the main README.md
- Read the documentation in `/docs`
- Review examples in `/examples`
- Open an issue on GitHub
