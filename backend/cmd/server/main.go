package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/api"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/kafka"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/middleware"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/services"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/ses"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize SES client
	sesClient, err := ses.NewSESClient()
	if err != nil {
		log.Fatal("Failed to initialize SES client:", err)
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer()
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer producer.Close()

	// Initialize services
	authService := services.NewAuthService(db)
	emailService := services.NewEmailService(db, sesClient, producer)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)
	emailHandler := api.NewEmailHandler(emailService)
	templateHandler := api.NewTemplateHandler(db)
	domainHandler := api.NewDomainHandler(db, sesClient)

	// Setup router
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(db))

	// Email routes
	protected.HandleFunc("/emails/send", emailHandler.SendEmail).Methods("POST")
	protected.HandleFunc("/emails/{id}", emailHandler.GetEmailStatus).Methods("GET")
	protected.HandleFunc("/emails/{id}/logs", emailHandler.GetEmailLogs).Methods("GET")
	protected.HandleFunc("/analytics", emailHandler.GetAnalytics).Methods("GET")
	protected.HandleFunc("/quota", emailHandler.GetQuota).Methods("GET")

	// Template routes
	protected.HandleFunc("/templates", templateHandler.CreateTemplate).Methods("POST")
	protected.HandleFunc("/templates", templateHandler.GetTemplates).Methods("GET")
	protected.HandleFunc("/templates/{id}", templateHandler.GetTemplate).Methods("GET")
	protected.HandleFunc("/templates/{id}", templateHandler.DeleteTemplate).Methods("DELETE")

	// Domain routes
	protected.HandleFunc("/domains", domainHandler.AddDomain).Methods("POST")
	protected.HandleFunc("/domains", domainHandler.GetDomains).Methods("GET")
	protected.HandleFunc("/domains/{id}/verify", domainHandler.VerifyDomain).Methods("POST")
	protected.HandleFunc("/domains/{id}", domainHandler.DeleteDomain).Methods("DELETE")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
