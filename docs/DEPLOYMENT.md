# Deployment Guide

This guide provides instructions for deploying UJUMBE to production.

## Prerequisites

- Server with Docker and Docker Compose
- Domain name configured
- AWS SES account with verified domain
- M-Pesa API credentials
- PostgreSQL database (can use Docker)
- Kafka cluster (can use Docker)

## Environment Configuration

1. **Copy and configure environment file**

```bash
cp .env.example .env
```

2. **Set required variables**

```bash
# Database
DATABASE_URL=postgresql://user:password@db-host:5432/ujumbe

# Kafka
KAFKA_BROKERS=kafka-host:9092

# AWS SES
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_SES_FROM_EMAIL=noreply@yourdomain.com

# M-Pesa (Production)
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_SHORTCODE=your_shortcode
MPESA_PASSKEY=your_passkey
MPESA_CALLBACK_URL=https://yourdomain.com/api/v1/mpesa/callback
MPESA_ENVIRONMENT=production

# Application
APP_ENV=production
APP_PORT=8080
JWT_SECRET=your_very_secure_random_secret_key
API_BASE_URL=https://yourdomain.com

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://dashboard.yourdomain.com
```

## Docker Deployment

### Option 1: Docker Compose (Recommended for small-medium deployments)

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Check status
docker-compose ps
```

### Option 2: Docker Swarm (For high availability)

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml ujumbe

# Check services
docker service ls

# Scale API service
docker service scale ujumbe_api=3
```

## Manual Deployment

### Backend

```bash
# Build
cd backend
go build -o ujumbe ./cmd/api

# Run with systemd
sudo cp ujumbe /usr/local/bin/
sudo cp ujumbe.service /etc/systemd/system/
sudo systemctl enable ujumbe
sudo systemctl start ujumbe
```

**ujumbe.service example:**
```ini
[Unit]
Description=UJUMBE Email API
After=network.target postgresql.service

[Service]
Type=simple
User=ujumbe
WorkingDirectory=/opt/ujumbe
EnvironmentFile=/opt/ujumbe/.env
ExecStart=/usr/local/bin/ujumbe
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Frontend

```bash
# Build
cd frontend
npm install
npm run build

# Serve with nginx
sudo cp -r dist/* /var/www/ujumbe/
```

**nginx configuration:**
```nginx
server {
    listen 80;
    server_name dashboard.yourdomain.com;
    
    root /var/www/ujumbe;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## SSL/TLS Setup

### Using Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com -d dashboard.yourdomain.com

# Auto-renewal is set up automatically
```

## Database Migration

```bash
# The application creates tables automatically on first run
# For migrations, use a tool like golang-migrate

# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir database/migrations -seq add_new_table

# Run migrations
migrate -path database/migrations -database "$DATABASE_URL" up
```

## Monitoring

### Logging

```bash
# View application logs
docker-compose logs -f api

# Or with journalctl
sudo journalctl -u ujumbe -f
```

### Health Check

```bash
curl http://localhost:8080/health
```

### Metrics

Consider adding:
- Prometheus for metrics
- Grafana for visualization
- Alertmanager for alerts

## Backup

### Database Backup

```bash
# Backup
pg_dump -h localhost -U ujumbe -d ujumbe > backup.sql

# Restore
psql -h localhost -U ujumbe -d ujumbe < backup.sql

# Automated daily backup
0 2 * * * pg_dump -h localhost -U ujumbe -d ujumbe > /backups/ujumbe-$(date +\%Y\%m\%d).sql
```

## Scaling

### Horizontal Scaling

1. **API Servers**: Run multiple instances behind a load balancer
2. **Kafka Consumers**: Run multiple consumer instances
3. **Database**: Use read replicas for read-heavy operations
4. **Kafka**: Increase partitions and brokers

### Load Balancer Configuration (nginx)

```nginx
upstream api_backend {
    least_conn;
    server api1:8080;
    server api2:8080;
    server api3:8080;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Security Checklist

- [ ] Use strong JWT secret
- [ ] Enable HTTPS
- [ ] Set secure CORS origins
- [ ] Use environment variables for secrets
- [ ] Enable database SSL
- [ ] Set up firewall rules
- [ ] Regular security updates
- [ ] Monitor for suspicious activity
- [ ] Implement rate limiting
- [ ] Regular backups

## Troubleshooting

### API won't start
- Check environment variables
- Verify database connection
- Check port availability
- Review logs

### Emails not sending
- Verify AWS SES credentials
- Check domain verification in AWS SES
- Review Kafka connection
- Check email logs in database

### M-Pesa not working
- Verify credentials
- Check callback URL is publicly accessible
- Ensure using correct environment (sandbox/production)
- Review transaction logs

## Performance Tuning

### PostgreSQL
```sql
-- Increase connection pool
max_connections = 200

-- Tune for better performance
shared_buffers = 256MB
effective_cache_size = 1GB
```

### Go API
```go
// Adjust server timeouts
ReadTimeout:  15 * time.Second,
WriteTimeout: 15 * time.Second,
IdleTimeout:  60 * time.Second,
```

### Kafka
```yaml
# Increase throughput
batch.size: 16384
linger.ms: 10
compression.type: snappy
```

## Maintenance

### Regular Tasks
- Monitor disk space
- Review logs weekly
- Update dependencies monthly
- Review security advisories
- Test backup restoration
- Monitor email delivery rates
- Check API response times

## Support

For deployment issues:
- Check documentation
- Review logs
- Open GitHub issue
- Contact support team

---

Happy deploying! 🚀
