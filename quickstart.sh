#!/bin/bash

# UJUMBE Quick Start Script
# This script helps you get UJUMBE up and running quickly

set -e

echo "🚀 UJUMBE Email Platform - Quick Start"
echo "======================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker and Docker Compose are installed${NC}"

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠ No .env file found. Creating from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Created .env file${NC}"
    echo -e "${YELLOW}⚠ Please edit .env with your credentials before proceeding${NC}"
    read -p "Press enter to continue after editing .env, or Ctrl+C to exit..."
fi

# Build and start services
echo ""
echo "📦 Building Docker images..."
docker-compose build

echo ""
echo "🚀 Starting services..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓ Services are running!${NC}"
else
    echo -e "${RED}❌ Some services failed to start. Check logs with: docker-compose logs${NC}"
    exit 1
fi

echo ""
echo "================================================"
echo -e "${GREEN}🎉 UJUMBE is now running!${NC}"
echo "================================================"
echo ""
echo "📍 API Server:     http://localhost:8080"
echo "📍 Dashboard:      http://localhost:5173 (if frontend is running)"
echo "📍 Health Check:   http://localhost:8080/health"
echo ""
echo "📝 Next steps:"
echo "   1. Visit http://localhost:8080/health to verify API is running"
echo "   2. Register a new user account"
echo "   3. Create an API key"
echo "   4. Send your first email!"
echo ""
echo "📚 Documentation:"
echo "   - API Examples: docs/API_EXAMPLES.md"
echo "   - Full Guide:   docs/README.md"
echo "   - Deployment:   docs/DEPLOYMENT.md"
echo ""
echo "🛠️  Useful commands:"
echo "   - View logs:        docker-compose logs -f"
echo "   - Stop services:    docker-compose down"
echo "   - Restart:          docker-compose restart"
echo ""
echo "💡 To run frontend development server:"
echo "   cd frontend && npm install && npm run dev"
echo ""
