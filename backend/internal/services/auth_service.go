package services

import (
	"database/sql"
	"fmt"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *database.DB
}

func NewAuthService(db *database.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	// Check if user exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return &models.RegisterResponse{
			Success: false,
			Message: "Email already registered",
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate API key
	apiKey := uuid.New().String()
	userID := uuid.New()

	// Create user
	_, err = s.db.Exec(`
		INSERT INTO users (id, email, password_hash, api_key)
		VALUES ($1, $2, $3, $4)
	`, userID, req.Email, string(hashedPassword), apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create initial quota
	_, err = s.db.Exec(`
		INSERT INTO user_quotas (user_id, free_emails_remaining, paid_emails_balance, monthly_reset_date)
		VALUES ($1, 100, 0, CURRENT_DATE)
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create quota: %w", err)
	}

	return &models.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		APIKey:  apiKey,
		UserID:  userID.String(),
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, api_key
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.APIKey)

	if err == sql.ErrNoRows {
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}

	return &models.LoginResponse{
		Success: true,
		Token:   user.APIKey,
		APIKey:  user.APIKey,
		Message: "Login successful",
	}, nil
}
