# Contributing to UJUMBE

Thank you for your interest in contributing to UJUMBE! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Prioritize the community's needs

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/Ismael-Njihia/UJUMBE/issues)
2. If not, create a new issue with:
   - Clear, descriptive title
   - Steps to reproduce
   - Expected vs actual behavior
   - System information (OS, Go version, etc.)
   - Relevant logs or error messages

### Suggesting Features

1. Check existing issues and discussions
2. Create a new issue with:
   - Clear description of the feature
   - Use cases and benefits
   - Potential implementation approach
   - Any alternatives considered

### Contributing Code

#### Setup Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/UJUMBE.git
   cd UJUMBE
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/Ismael-Njihia/UJUMBE.git
   ```

4. Set up development environment:
   ```bash
   # Backend
   cd backend
   go mod download
   
   # Frontend
   cd ../frontend
   npm install
   ```

5. Start services:
   ```bash
   make run
   ```

#### Making Changes

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following our coding standards

3. Test your changes:
   ```bash
   # Backend tests
   cd backend
   go test ./...
   
   # Build to ensure compilation
   go build ./cmd/server
   go build ./cmd/worker
   ```

4. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Use conventional commit messages:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes (formatting, etc.)
   - `refactor:` - Code refactoring
   - `test:` - Adding or updating tests
   - `chore:` - Maintenance tasks

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a Pull Request from your fork to the main repository

#### Pull Request Guidelines

- Link related issues in the PR description
- Provide clear description of changes
- Include screenshots for UI changes
- Ensure all tests pass
- Keep PRs focused on a single feature/fix
- Update documentation if needed
- Add tests for new features

### Coding Standards

#### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` and fix issues
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors explicitly

Example:
```go
// SendEmail sends an email via AWS SES and tracks it in the database.
// It returns the email ID and any error encountered.
func (s *EmailService) SendEmail(userID uuid.UUID, req models.SendEmailRequest) (*models.SendEmailResponse, error) {
    // Validate input
    if req.To == "" {
        return nil, errors.New("recipient email is required")
    }
    
    // Implementation
    // ...
    
    return response, nil
}
```

#### JavaScript/Svelte Code

- Use 2 spaces for indentation
- Use meaningful variable names
- Add JSDoc comments for complex functions
- Keep components small and reusable
- Use semantic HTML

#### Database Migrations

- Create sequential migration files
- Include both up and down migrations
- Test migrations on sample data
- Document schema changes

#### API Design

- Follow RESTful conventions
- Use proper HTTP status codes
- Include error messages in responses
- Version APIs appropriately
- Document new endpoints

### Testing

#### Backend Tests

```go
func TestSendEmail(t *testing.T) {
    // Setup
    mockDB := setupMockDB(t)
    service := NewEmailService(mockDB, mockSES, mockKafka)
    
    // Test
    response, err := service.SendEmail(testUserID, testRequest)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, response.EmailID)
}
```

#### Integration Tests

- Test complete workflows
- Use test database
- Clean up after tests
- Mock external services when appropriate

### Documentation

- Update README.md for major changes
- Add API documentation in `docs/API.md`
- Update setup guide if needed
- Include code examples
- Keep tutorials up to date

### Areas Needing Contribution

We especially welcome contributions in these areas:

1. **M-Pesa Integration**
   - Complete M-Pesa STK Push implementation
   - Callback handling
   - Transaction verification

2. **Testing**
   - Unit tests for services
   - Integration tests
   - End-to-end tests
   - Load testing

3. **Features**
   - Email scheduling
   - Attachment support
   - Email preview
   - Webhook notifications
   - Multi-language support

4. **Performance**
   - Database query optimization
   - Caching strategies
   - Connection pooling
   - Rate limiting

5. **Security**
   - Security audit
   - Penetration testing
   - Rate limiting per user
   - API throttling

6. **Documentation**
   - Video tutorials
   - More code examples
   - Translation to Swahili
   - Architecture diagrams

7. **DevOps**
   - Kubernetes deployment
   - CI/CD pipelines
   - Monitoring dashboards
   - Backup strategies

### Development Tips

#### Running Locally

```bash
# Start all services
make run

# View logs
make logs

# Stop services
make stop

# Clean up
make clean
```

#### Database Access

```bash
# Connect to database
docker exec -it ujumbe-postgres psql -U postgres -d ujumbe

# Run migrations
make migrate-up
```

#### Debugging

```bash
# Backend logs
docker logs -f ujumbe-backend

# Worker logs
docker logs -f ujumbe-worker

# Kafka consumer
docker exec -it ujumbe-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic email_jobs
```

### Git Workflow

1. Keep your fork up to date:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. Rebase your feature branch:
   ```bash
   git checkout feature/your-feature
   git rebase main
   ```

3. Resolve conflicts if any

4. Force push if needed (only on your branch):
   ```bash
   git push -f origin feature/your-feature
   ```

### Review Process

1. Automated checks run on your PR
2. Maintainers review your code
3. Address any feedback
4. Once approved, PR will be merged

### Getting Help

- Ask questions in GitHub Discussions
- Join our community channels
- Check existing documentation
- Review closed issues for similar problems

### Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Given credit in documentation

## Project Structure

```
UJUMBE/
├── backend/           # Go backend
│   ├── cmd/          # Application entry points
│   ├── internal/     # Internal packages
│   └── pkg/          # Public packages
├── frontend/         # Svelte frontend
│   └── src/         # Source files
├── database/        # Database migrations
├── docs/           # Documentation
├── examples/       # Integration examples
└── docker-compose.yml
```

## License

By contributing to UJUMBE, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to:
- Open an issue for questions
- Start a discussion
- Email: support@ujumbe.co.ke

Thank you for contributing to UJUMBE! 🎉
