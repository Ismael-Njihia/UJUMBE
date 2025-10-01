package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/auth"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/billing"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/domain"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/email"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/kafka"
	"github.com/Ismael-Njihia/UJUMBE/backend/pkg/config"
	"github.com/Ismael-Njihia/UJUMBE/backend/pkg/mpesa"
	"github.com/Ismael-Njihia/UJUMBE/backend/pkg/ses"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type API struct {
	authService    *auth.Service
	emailService   *email.Service
	domainService  *domain.Service
	billingService *billing.Service
	kafkaProducer  *kafka.Producer
	sesClient      *ses.SESClient
	mpesaClient    *mpesa.Client
	config         *config.Config
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize services
	authService := auth.NewService(db.DB, cfg.JWTSecret)
	emailService := email.NewService(db.DB)
	domainService := domain.NewService(db.DB)
	billingService := billing.NewService(db.DB)

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaEmailTopic)
	defer kafkaProducer.Close()

	// Initialize SES client
	sesClient, err := ses.NewSESClient(cfg.AWSRegion, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey)
	if err != nil {
		log.Printf("Warning: Failed to initialize SES client: %v", err)
	}

	// Initialize M-Pesa client
	mpesaClient := mpesa.NewClient(
		cfg.MpesaConsumerKey,
		cfg.MpesaConsumerSecret,
		cfg.MpesaShortcode,
		cfg.MpesaPasskey,
		cfg.MpesaCallbackURL,
		cfg.MpesaEnvironment,
	)

	// Create API instance
	api := &API{
		authService:    authService,
		emailService:   emailService,
		domainService:  domainService,
		billingService: billingService,
		kafkaProducer:  kafkaProducer,
		sesClient:      sesClient,
		mpesaClient:    mpesaClient,
		config:         cfg,
	}

	// Setup router
	router := mux.NewRouter()
	api.setupRoutes(router)

	// Setup CORS
	origins := strings.Split(cfg.CORSAllowedOrigins, ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Setup server
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Starting server on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start Kafka consumer in a separate goroutine
	go api.startEmailConsumer(cfg)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func (api *API) setupRoutes(router *mux.Router) {
	// Health check
	router.HandleFunc("/health", api.healthCheck).Methods("GET")

	// API v1
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	v1.HandleFunc("/auth/register", api.register).Methods("POST")
	v1.HandleFunc("/auth/login", api.login).Methods("POST")

	// Protected routes
	protected := v1.PathPrefix("").Subrouter()
	protected.Use(api.authMiddleware)

	// User routes
	protected.HandleFunc("/user", api.getUser).Methods("GET")
	protected.HandleFunc("/user/stats", api.getDashboardStats).Methods("GET")

	// API Key routes
	protected.HandleFunc("/api-keys", api.getAPIKeys).Methods("GET")
	protected.HandleFunc("/api-keys", api.createAPIKey).Methods("POST")
	protected.HandleFunc("/api-keys/{id}", api.revokeAPIKey).Methods("DELETE")

	// Domain routes
	protected.HandleFunc("/domains", api.getDomains).Methods("GET")
	protected.HandleFunc("/domains", api.addDomain).Methods("POST")
	protected.HandleFunc("/domains/{id}", api.getDomain).Methods("GET")
	protected.HandleFunc("/domains/{id}/verify", api.verifyDomain).Methods("POST")
	protected.HandleFunc("/domains/{id}", api.deleteDomain).Methods("DELETE")

	// Template routes
	protected.HandleFunc("/templates", api.getTemplates).Methods("GET")
	protected.HandleFunc("/templates", api.createTemplate).Methods("POST")
	protected.HandleFunc("/templates/{id}", api.getTemplate).Methods("GET")
	protected.HandleFunc("/templates/{id}", api.updateTemplate).Methods("PUT")
	protected.HandleFunc("/templates/{id}", api.deleteTemplate).Methods("DELETE")

	// Email routes (can use API key or JWT)
	v1.HandleFunc("/emails/send", api.authOrAPIKeyMiddleware(api.sendEmail)).Methods("POST")

	// Email logs routes
	protected.HandleFunc("/emails/logs", api.getEmailLogs).Methods("GET")

	// Analytics routes
	protected.HandleFunc("/analytics", api.getAnalytics).Methods("GET")

	// Billing routes
	protected.HandleFunc("/billing/transactions", api.getTransactions).Methods("GET")
	protected.HandleFunc("/billing/topup", api.initiateTopup).Methods("POST")

	// M-Pesa callback (public)
	v1.HandleFunc("/mpesa/callback", api.mpesaCallback).Methods("POST")
}

func (api *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Middleware to authenticate JWT token
func (api *API) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := api.authService.ValidateToken(tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to authenticate either JWT or API key
func (api *API) authOrAPIKeyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try JWT first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := api.authService.ValidateToken(tokenString)
			if err == nil {
				ctx := context.WithValue(r.Context(), "userID", userID)
				handler(w, r.WithContext(ctx))
				return
			}
		}

		// Try API key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			userID, err := api.authService.ValidateAPIKey(apiKey)
			if err == nil {
				ctx := context.WithValue(r.Context(), "userID", userID)
				handler(w, r.WithContext(ctx))
				return
			}
		}

		respondError(w, http.StatusUnauthorized, "Authentication required")
	}
}

func (api *API) startEmailConsumer(cfg *config.Config) {
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaEmailTopic, cfg.KafkaGroupID)
	defer consumer.Close()

	ctx := context.Background()
	err := consumer.ConsumeEmails(ctx, func(job *kafka.EmailJob) error {
		// Send email via SES
		if api.sesClient == nil {
			log.Printf("SES client not initialized, skipping email: %s", job.ID)
			return nil
		}

		messageID, err := api.sesClient.SendEmail(&ses.EmailMessage{
			From:     job.From,
			To:       []string{job.To},
			Subject:  job.Subject,
			HTMLBody: job.HTMLBody,
			TextBody: job.TextBody,
		})

		status := "sent"
		errorMessage := ""
		if err != nil {
			status = "failed"
			errorMessage = err.Error()
			log.Printf("Failed to send email: %v", err)
		}

		// Update log
		if err := api.emailService.UpdateEmailLogStatus(job.ID, status, errorMessage, messageID); err != nil {
			log.Printf("Failed to update email log: %v", err)
		}

		// Record analytics
		emailsSent := 0
		failed := 0
		if status == "sent" {
			emailsSent = 1
		} else {
			failed = 1
		}
		if err := api.emailService.RecordAnalytics(job.UserID, emailsSent, failed, 0); err != nil {
			log.Printf("Failed to record analytics: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Email consumer error: %v", err)
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func getUserID(r *http.Request) string {
	userID, _ := r.Context().Value("userID").(string)
	return userID
}
