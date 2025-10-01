package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURL  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string

	// Kafka
	KafkaBrokers   string
	KafkaEmailTopic string
	KafkaGroupID   string

	// AWS SES
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSESFromEmail    string

	// M-Pesa
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	MpesaShortcode      string
	MpesaPasskey        string
	MpesaCallbackURL    string
	MpesaEnvironment    string

	// Application
	AppEnv         string
	AppPort        string
	JWTSecret      string
	APIBaseURL     string

	// Email Quotas
	FreeEmailQuota    int
	EmailPricePerUnit float64

	// CORS
	CORSAllowedOrigins string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	freeQuota, _ := strconv.Atoi(getEnv("FREE_EMAIL_QUOTA", "100"))
	pricePerUnit, _ := strconv.ParseFloat(getEnv("EMAIL_PRICE_PER_UNIT", "1.0"), 64)

	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgresql://ujumbe:password@localhost:5432/ujumbe?sslmode=disable"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "ujumbe"),
		DBPassword:          getEnv("DB_PASSWORD", "password"),
		DBName:              getEnv("DB_NAME", "ujumbe"),
		KafkaBrokers:        getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaEmailTopic:     getEnv("KAFKA_EMAIL_TOPIC", "emails"),
		KafkaGroupID:        getEnv("KAFKA_GROUP_ID", "ujumbe-consumer"),
		AWSRegion:           getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:      getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSSESFromEmail:     getEnv("AWS_SES_FROM_EMAIL", "noreply@yourdomain.com"),
		MpesaConsumerKey:    getEnv("MPESA_CONSUMER_KEY", ""),
		MpesaConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
		MpesaShortcode:      getEnv("MPESA_SHORTCODE", "174379"),
		MpesaPasskey:        getEnv("MPESA_PASSKEY", ""),
		MpesaCallbackURL:    getEnv("MPESA_CALLBACK_URL", ""),
		MpesaEnvironment:    getEnv("MPESA_ENVIRONMENT", "sandbox"),
		AppEnv:              getEnv("APP_ENV", "development"),
		AppPort:             getEnv("APP_PORT", "8080"),
		JWTSecret:           getEnv("JWT_SECRET", "change-this-secret-key"),
		APIBaseURL:          getEnv("API_BASE_URL", "http://localhost:8080"),
		FreeEmailQuota:      freeQuota,
		EmailPricePerUnit:   pricePerUnit,
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
