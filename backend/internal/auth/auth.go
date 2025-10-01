package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidAPIKey      = errors.New("invalid API key")
)

type Service struct {
	db        *sql.DB
	jwtSecret string
}

func NewService(db *sql.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user
func (s *Service) Register(email, password, name string) (*models.User, error) {
	// Check if user exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		EmailQuota:   100,
		EmailsSent:   0,
		Balance:      0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.db.Exec(
		`INSERT INTO users (id, email, password_hash, name, email_quota, emails_sent, balance, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.PasswordHash, user.Name, user.EmailQuota, user.EmailsSent, user.Balance, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(email, password string) (string, *models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`SELECT id, email, password_hash, name, email_quota, emails_sent, balance, created_at, updated_at 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.EmailQuota, &user.EmailsSent, &user.Balance, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return "", nil, ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &user, nil
}

// ValidateToken validates a JWT token
func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", errors.New("invalid token")
}

// CreateAPIKey generates a new API key for a user
func (s *Service) CreateAPIKey(userID, name string) (*models.APIKey, error) {
	// Generate random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}
	key := "ujumbe_" + base64.URLEncoding.EncodeToString(keyBytes)

	apiKey := &models.APIKey{
		ID:        uuid.New().String(),
		UserID:    userID,
		Key:       key,
		Name:      name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Exec(
		`INSERT INTO api_keys (id, user_id, key, name, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		apiKey.ID, apiKey.UserID, apiKey.Key, apiKey.Name, apiKey.IsActive, apiKey.CreatedAt, apiKey.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

// ValidateAPIKey validates an API key and returns the user ID
func (s *Service) ValidateAPIKey(key string) (string, error) {
	var userID string
	var isActive bool
	err := s.db.QueryRow(
		`SELECT user_id, is_active FROM api_keys WHERE key = $1`,
		key,
	).Scan(&userID, &isActive)

	if err == sql.ErrNoRows || !isActive {
		return "", ErrInvalidAPIKey
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetUserAPIKeys returns all API keys for a user
func (s *Service) GetUserAPIKeys(userID string) ([]models.APIKey, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, key, name, is_active, created_at, updated_at 
		 FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		if err := rows.Scan(&apiKey.ID, &apiKey.UserID, &apiKey.Key, &apiKey.Name, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.UpdatedAt); err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// RevokeAPIKey deactivates an API key
func (s *Service) RevokeAPIKey(keyID, userID string) error {
	_, err := s.db.Exec(
		`UPDATE api_keys SET is_active = FALSE, updated_at = $1 WHERE id = $2 AND user_id = $3`,
		time.Now(), keyID, userID,
	)
	return err
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`SELECT id, email, name, email_quota, emails_sent, balance, created_at, updated_at 
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.EmailQuota, &user.EmailsSent, &user.Balance, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
