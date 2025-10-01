# UJUMBE Setup Guide

This guide will walk you through setting up UJUMBE for development and production.

## Prerequisites

Before you begin, ensure you have the following installed:
- Docker (v20.10+)
- Docker Compose (v2.0+)
- Go (v1.21+) - for local development
- Node.js (v18+) - for frontend development
- Make - for using Makefile commands
- AWS Account with SES access

## AWS SES Setup

### 1. Create AWS Account and Configure SES

1. Sign up for AWS: https://aws.amazon.com
2. Navigate to SES (Simple Email Service)
3. Request production access (if not already granted)
4. Verify your sender email addresses or domains

### 2. Get AWS Credentials

1. Go to AWS IAM Console
2. Create a new IAM user for UJUMBE
3. Attach the `AmazonSESFullAccess` policy
4. Generate access keys
5. Save the Access Key ID and Secret Access Key

### 3. Verify Email Addresses (Development)

For development, verify individual email addresses:
1. In SES Console, go to "Verified identities"
2. Click "Create identity"
3. Select "Email address"
4. Enter your email and click "Create identity"
5. Check your email and click the verification link

### 4. Verify Domains (Production)

For production, verify entire domains:
1. In SES Console, go to "Verified identities"
2. Click "Create identity"
3. Select "Domain"
4. Enter your domain name
5. Add the provided DNS records to your domain registrar
6. Wait for verification (usually 5-10 minutes)

## Database Setup

### Using Docker (Recommended)

The PostgreSQL database is included in docker-compose.yml:
```bash
make run
```

### Manual Setup

If you prefer to set up PostgreSQL manually:

1. Install PostgreSQL:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install postgresql-15

   # macOS
   brew install postgresql@15
   ```

2. Create database and user:
   ```sql
   CREATE DATABASE ujumbe;
   CREATE USER ujumbe_user WITH PASSWORD 'secure_password';
   GRANT ALL PRIVILEGES ON DATABASE ujumbe TO ujumbe_user;
   ```

3. Run migrations:
   ```bash
   psql -U ujumbe_user -d ujumbe -f database/migrations/001_init_schema.sql
   ```

## Kafka Setup

### Using Docker (Recommended)

Kafka and Zookeeper are included in docker-compose.yml:
```bash
make run
```

### Manual Setup

1. Download Kafka: https://kafka.apache.org/downloads
2. Extract and navigate to the Kafka directory
3. Start Zookeeper:
   ```bash
   bin/zookeeper-server-start.sh config/zookeeper.properties
   ```
4. Start Kafka:
   ```bash
   bin/kafka-server-start.sh config/server.properties
   ```
5. Create topic:
   ```bash
   bin/kafka-topics.sh --create --topic email_jobs --bootstrap-server localhost:9092
   ```

## Backend Setup

### 1. Configure Environment Variables

```bash
cd backend
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ujumbe
DB_SSLMODE=disable

# AWS SES
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_actual_access_key
AWS_SECRET_ACCESS_KEY=your_actual_secret_key
AWS_SES_FROM_EMAIL=noreply@yourdomain.com

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=email_jobs
KAFKA_GROUP_ID=ujumbe_consumers

# Server
SERVER_PORT=8080
JWT_SECRET=generate_a_random_secret_key_here

# Application
FREE_EMAIL_QUOTA=100
EMAIL_PRICE_PER_UNIT=1.0
ENVIRONMENT=development
```

### 2. Install Dependencies

```bash
cd backend
go mod download
```

### 3. Run the API Server

Using Make:
```bash
make dev
```

Or directly:
```bash
cd backend
go run cmd/server/main.go
```

### 4. Run the Worker

In a new terminal:
```bash
make worker-dev
```

Or directly:
```bash
cd backend
go run cmd/worker/main.go
```

## Frontend Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure API Endpoint

The API endpoint is configured in `frontend/src/lib/api.js`:
```javascript
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

For production, update this to your production API URL.

### 3. Run Development Server

```bash
npm run dev
```

Access the dashboard at: http://localhost:3000

### 4. Build for Production

```bash
npm run build
```

Built files will be in `frontend/dist/`

## Docker Setup (Recommended for Production)

### 1. Configure Environment

Create a `.env` file in the project root with production values.

### 2. Build and Run

```bash
# Build images
make build

# Start all services
make run

# View logs
make logs

# Stop services
make stop
```

### 3. Verify Services

Check if all services are running:
```bash
docker ps
```

You should see:
- ujumbe-postgres
- ujumbe-kafka
- ujumbe-zookeeper
- ujumbe-backend
- ujumbe-worker

### 4. Health Check

```bash
curl http://localhost:8080/health
```

Should return:
```json
{"status":"healthy"}
```

## Testing the Setup

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }'
```

Save the returned `api_key`.

### 2. Verify a Domain

```bash
curl -X POST http://localhost:8080/api/v1/domains \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "domain": "yourdomain.com"
  }'
```

Follow AWS SES instructions to add DNS records.

### 3. Send a Test Email

```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "from": "noreply@yourdomain.com",
    "to": "recipient@example.com",
    "subject": "Test Email from UJUMBE",
    "html_body": "<h1>Hello!</h1><p>This is a test email.</p>"
  }'
```

## M-Pesa Integration (Optional)

To enable M-Pesa payments:

### 1. Get M-Pesa Credentials

1. Register on Safaricom Daraja: https://developer.safaricom.co.ke
2. Create an app
3. Get Consumer Key and Consumer Secret
4. Note your shortcode and passkey

### 2. Configure M-Pesa

Add to `.env`:
```env
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_SHORTCODE=174379
MPESA_PASSKEY=your_passkey
MPESA_CALLBACK_URL=https://yourdomain.com/api/v1/mpesa/callback
```

### 3. Implement M-Pesa Handler

Create `backend/internal/api/mpesa_handler.go` with M-Pesa STK push and callback logic.

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs ujumbe-postgres

# Connect to database
docker exec -it ujumbe-postgres psql -U postgres -d ujumbe
```

### Kafka Connection Issues

```bash
# Check if Kafka is running
docker ps | grep kafka

# Check logs
docker logs ujumbe-kafka

# List topics
docker exec -it ujumbe-kafka kafka-topics --list --bootstrap-server localhost:9092
```

### AWS SES Issues

- Ensure your AWS credentials are correct
- Check if you're in sandbox mode (can only send to verified emails)
- Verify sender domain/email in SES console
- Check SES sending limits

### Email Not Sending

1. Check worker logs:
   ```bash
   docker logs ujumbe-worker
   ```

2. Check Kafka for pending jobs:
   ```bash
   docker exec -it ujumbe-kafka kafka-console-consumer \
     --bootstrap-server localhost:9092 \
     --topic email_jobs \
     --from-beginning
   ```

3. Check email status via API:
   ```bash
   curl http://localhost:8080/api/v1/emails/{email_id} \
     -H "X-API-Key: your-api-key"
   ```

## Production Considerations

### Security
- Use strong passwords and secrets
- Enable SSL/TLS for all connections
- Use environment-specific `.env` files
- Implement rate limiting
- Set up firewall rules
- Regular security audits

### Performance
- Scale workers horizontally for high volume
- Use connection pooling for database
- Configure Kafka partitions for throughput
- Use CDN for frontend assets
- Implement caching where appropriate

### Monitoring
- Set up application logging
- Use Prometheus for metrics
- Configure alerting for failures
- Monitor AWS SES sending statistics
- Track Kafka consumer lag

### Backup
- Regular database backups
- Back up configuration files
- Document disaster recovery procedures

## Getting Help

- Check logs: `make logs`
- Review API documentation: [docs/API.md](API.md)
- Check GitHub issues: https://github.com/Ismael-Njihia/UJUMBE/issues
- Contact support: support@ujumbe.co.ke

## Next Steps

After setup:
1. Explore the API documentation
2. Create email templates
3. Verify your domains
4. Integrate with your application
5. Monitor email delivery
6. Set up M-Pesa billing (optional)

Happy emailing! 📧
