package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func New(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) InitSchema() error {
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		email_quota INTEGER DEFAULT 100,
		emails_sent INTEGER DEFAULT 0,
		balance DECIMAL(10, 2) DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- API Keys table
	CREATE TABLE IF NOT EXISTS api_keys (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		key VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Domains table
	CREATE TABLE IF NOT EXISTS domains (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		domain VARCHAR(255) NOT NULL,
		is_verified BOOLEAN DEFAULT FALSE,
		verification_code VARCHAR(255),
		dkim_status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, domain)
	);

	-- Email Templates table
	CREATE TABLE IF NOT EXISTS email_templates (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		subject VARCHAR(500) NOT NULL,
		html_body TEXT NOT NULL,
		text_body TEXT,
		variables JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Email Logs table
	CREATE TABLE IF NOT EXISTS email_logs (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		template_id VARCHAR(36) REFERENCES email_templates(id) ON DELETE SET NULL,
		from_email VARCHAR(255) NOT NULL,
		to_email VARCHAR(255) NOT NULL,
		subject VARCHAR(500) NOT NULL,
		status VARCHAR(50) NOT NULL,
		error_message TEXT,
		message_id VARCHAR(255),
		sent_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Transactions table
	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		amount DECIMAL(10, 2) NOT NULL,
		currency VARCHAR(10) DEFAULT 'KES',
		type VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL,
		mpesa_receipt_no VARCHAR(255),
		phone_number VARCHAR(20),
		emails_purchased INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Analytics table
	CREATE TABLE IF NOT EXISTS analytics (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		date DATE NOT NULL,
		emails_sent INTEGER DEFAULT 0,
		failed INTEGER DEFAULT 0,
		bounced INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, date)
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
	CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
	CREATE INDEX IF NOT EXISTS idx_domains_user_id ON domains(user_id);
	CREATE INDEX IF NOT EXISTS idx_email_templates_user_id ON email_templates(user_id);
	CREATE INDEX IF NOT EXISTS idx_email_logs_user_id ON email_logs(user_id);
	CREATE INDEX IF NOT EXISTS idx_email_logs_status ON email_logs(status);
	CREATE INDEX IF NOT EXISTS idx_email_logs_created_at ON email_logs(created_at);
	CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
	CREATE INDEX IF NOT EXISTS idx_analytics_user_id ON analytics(user_id);
	CREATE INDEX IF NOT EXISTS idx_analytics_date ON analytics(date);
	`

	_, err := d.DB.Exec(schema)
	return err
}
