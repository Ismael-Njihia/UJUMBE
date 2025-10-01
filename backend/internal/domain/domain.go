package domain

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrDomainExists   = errors.New("domain already exists")
	ErrDomainNotFound = errors.New("domain not found")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// AddDomain adds a new domain for verification
func (s *Service) AddDomain(userID, domainName string) (*models.Domain, error) {
	// Check if domain already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM domains WHERE user_id = $1 AND domain = $2)", userID, domainName).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDomainExists
	}

	// Generate verification code
	verificationCode, err := generateVerificationCode()
	if err != nil {
		return nil, err
	}

	domain := &models.Domain{
		ID:               uuid.New().String(),
		UserID:           userID,
		Domain:           domainName,
		IsVerified:       false,
		VerificationCode: verificationCode,
		DKIMStatus:       "pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = s.db.Exec(
		`INSERT INTO domains (id, user_id, domain, is_verified, verification_code, dkim_status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		domain.ID, domain.UserID, domain.Domain, domain.IsVerified, domain.VerificationCode, domain.DKIMStatus, domain.CreatedAt, domain.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

// GetDomain retrieves a domain by ID
func (s *Service) GetDomain(domainID, userID string) (*models.Domain, error) {
	var domain models.Domain
	err := s.db.QueryRow(
		`SELECT id, user_id, domain, is_verified, verification_code, dkim_status, created_at, updated_at
		 FROM domains WHERE id = $1 AND user_id = $2`,
		domainID, userID,
	).Scan(&domain.ID, &domain.UserID, &domain.Domain, &domain.IsVerified, &domain.VerificationCode, &domain.DKIMStatus, &domain.CreatedAt, &domain.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrDomainNotFound
	}
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

// GetUserDomains retrieves all domains for a user
func (s *Service) GetUserDomains(userID string) ([]models.Domain, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, domain, is_verified, verification_code, dkim_status, created_at, updated_at
		 FROM domains WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		if err := rows.Scan(&domain.ID, &domain.UserID, &domain.Domain, &domain.IsVerified, &domain.VerificationCode, &domain.DKIMStatus, &domain.CreatedAt, &domain.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

// VerifyDomain marks a domain as verified
func (s *Service) VerifyDomain(domainID, userID string) error {
	_, err := s.db.Exec(
		`UPDATE domains SET is_verified = TRUE, updated_at = $1 WHERE id = $2 AND user_id = $3`,
		time.Now(), domainID, userID,
	)
	return err
}

// UpdateDKIMStatus updates the DKIM status of a domain
func (s *Service) UpdateDKIMStatus(domainID, userID, status string) error {
	_, err := s.db.Exec(
		`UPDATE domains SET dkim_status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`,
		status, time.Now(), domainID, userID,
	)
	return err
}

// DeleteDomain deletes a domain
func (s *Service) DeleteDomain(domainID, userID string) error {
	_, err := s.db.Exec(
		`DELETE FROM domains WHERE id = $1 AND user_id = $2`,
		domainID, userID,
	)
	return err
}

// GetVerifiedDomains returns all verified domains for a user
func (s *Service) GetVerifiedDomains(userID string) ([]models.Domain, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, domain, is_verified, verification_code, dkim_status, created_at, updated_at
		 FROM domains WHERE user_id = $1 AND is_verified = TRUE ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		if err := rows.Scan(&domain.ID, &domain.UserID, &domain.Domain, &domain.IsVerified, &domain.VerificationCode, &domain.DKIMStatus, &domain.CreatedAt, &domain.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func generateVerificationCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
