# Contributing to UJUMBE

Thank you for your interest in contributing to UJUMBE! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/UJUMBE.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit with clear messages
7. Push to your fork
8. Create a Pull Request

## Development Setup

### Backend Development

```bash
# Install Go 1.21+
# Install dependencies
go mod download

# Run the server
go run backend/cmd/api/*.go

# Or use Make
make run
```

### Frontend Development

```bash
# Install Node.js 18+
cd frontend
npm install

# Run dev server
npm run dev
```

### Running with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Code Standards

### Go Code
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names
- Add comments for exported functions
- Write tests for new features
- Keep functions focused and small

### JavaScript/Svelte Code
- Use ES6+ features
- Follow Svelte best practices
- Use meaningful component names
- Keep components focused

### Commits
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Keep commits focused on a single change
- Reference issues when applicable

Example:
```
Add email template validation
Fix M-Pesa callback parsing issue
Update README with new API endpoints
```

## Testing

### Backend Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./backend/internal/auth
```

### Frontend Tests
```bash
cd frontend
npm test
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Request review from maintainers

## Areas for Contribution

- 🐛 Bug fixes
- ✨ New features
- 📝 Documentation improvements
- 🎨 UI/UX enhancements
- 🧪 Test coverage
- ♿ Accessibility improvements
- 🌍 Internationalization

## Questions?

Feel free to open an issue for:
- Questions about the codebase
- Feature requests
- Bug reports
- General discussion

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn

Thank you for contributing to UJUMBE! 🚀
