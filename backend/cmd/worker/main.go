package main

import (
	"log"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/kafka"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/services"
	"github.com/Ismael-Njihia/UJUMBE/backend/internal/ses"
	"github.com/joho/godotenv"
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

	// Initialize Kafka producer (needed for service)
	producer, err := kafka.NewProducer()
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer producer.Close()

	// Initialize email service
	emailService := services.NewEmailService(db, sesClient, producer)

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer()
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer:", err)
	}
	defer consumer.Close()

	log.Println("Email worker started. Waiting for jobs...")

	// Start consuming jobs
	err = consumer.ConsumeEmailJobs(func(job kafka.EmailJob) error {
		log.Printf("Processing email job: %s\n", job.EmailID)
		return emailService.ProcessEmailJob(job)
	})

	if err != nil {
		log.Fatal("Consumer error:", err)
	}
}
