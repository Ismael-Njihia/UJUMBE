package billing

import (
	"database/sql"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/models"
	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateTransaction creates a new billing transaction
func (s *Service) CreateTransaction(userID string, amount float64, transactionType, phoneNumber string, emailsPurchased int) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:              uuid.New().String(),
		UserID:          userID,
		Amount:          amount,
		Currency:        "KES",
		Type:            transactionType,
		Status:          "pending",
		PhoneNumber:     phoneNumber,
		EmailsPurchased: emailsPurchased,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err := s.db.Exec(
		`INSERT INTO transactions (id, user_id, amount, currency, type, status, phone_number, emails_purchased, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		transaction.ID, transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Status,
		transaction.PhoneNumber, transaction.EmailsPurchased, transaction.CreatedAt, transaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (s *Service) UpdateTransactionStatus(transactionID, status, mpesaReceiptNo string) error {
	_, err := s.db.Exec(
		`UPDATE transactions SET status = $1, mpesa_receipt_no = $2, updated_at = $3 WHERE id = $4`,
		status, mpesaReceiptNo, time.Now(), transactionID,
	)
	return err
}

// GetTransaction retrieves a transaction by ID
func (s *Service) GetTransaction(transactionID string) (*models.Transaction, error) {
	var transaction models.Transaction
	var mpesaReceiptNo, phoneNumber sql.NullString

	err := s.db.QueryRow(
		`SELECT id, user_id, amount, currency, type, status, mpesa_receipt_no, phone_number, emails_purchased, created_at, updated_at
		 FROM transactions WHERE id = $1`,
		transactionID,
	).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Currency, &transaction.Type,
		&transaction.Status, &mpesaReceiptNo, &phoneNumber, &transaction.EmailsPurchased, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if mpesaReceiptNo.Valid {
		transaction.MpesaReceiptNo = mpesaReceiptNo.String
	}
	if phoneNumber.Valid {
		transaction.PhoneNumber = phoneNumber.String
	}

	return &transaction, nil
}

// GetUserTransactions retrieves all transactions for a user
func (s *Service) GetUserTransactions(userID string, limit, offset int) ([]models.Transaction, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, amount, currency, type, status, mpesa_receipt_no, phone_number, emails_purchased, created_at, updated_at
		 FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var mpesaReceiptNo, phoneNumber sql.NullString

		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Currency, &transaction.Type,
			&transaction.Status, &mpesaReceiptNo, &phoneNumber, &transaction.EmailsPurchased, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}

		if mpesaReceiptNo.Valid {
			transaction.MpesaReceiptNo = mpesaReceiptNo.String
		}
		if phoneNumber.Valid {
			transaction.PhoneNumber = phoneNumber.String
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// AddBalance adds balance to a user's account
func (s *Service) AddBalance(userID string, amount float64) error {
	_, err := s.db.Exec(
		`UPDATE users SET balance = balance + $1, updated_at = $2 WHERE id = $3`,
		amount, time.Now(), userID,
	)
	return err
}

// AddEmailQuota adds email quota to a user's account
func (s *Service) AddEmailQuota(userID string, quota int) error {
	_, err := s.db.Exec(
		`UPDATE users SET email_quota = email_quota + $1, updated_at = $2 WHERE id = $3`,
		quota, time.Now(), userID,
	)
	return err
}

// DeductBalance deducts balance from a user's account
func (s *Service) DeductBalance(userID string, amount float64) error {
	_, err := s.db.Exec(
		`UPDATE users SET balance = balance - $1, updated_at = $2 WHERE id = $3`,
		amount, time.Now(), userID,
	)
	return err
}

// GetUserBalance retrieves a user's current balance
func (s *Service) GetUserBalance(userID string) (float64, error) {
	var balance float64
	err := s.db.QueryRow(
		`SELECT balance FROM users WHERE id = $1`,
		userID,
	).Scan(&balance)
	return balance, err
}

// ResetMonthlyQuota resets the monthly email quota for all users
func (s *Service) ResetMonthlyQuota(freeQuota int) error {
	_, err := s.db.Exec(
		`UPDATE users SET emails_sent = 0, email_quota = $1, updated_at = $2`,
		freeQuota, time.Now(),
	)
	return err
}
