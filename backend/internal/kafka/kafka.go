package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type EmailJob struct {
	EmailID   uuid.UUID              `json:"email_id"`
	UserID    uuid.UUID              `json:"user_id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Subject   string                 `json:"subject"`
	HTMLBody  string                 `json:"html_body"`
	TextBody  string                 `json:"text_body"`
}

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer() (*Producer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) ProduceEmailJob(job EmailJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report
	e := <-p.producer.Events()
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer() (*Consumer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	topic := os.Getenv("KAFKA_TOPIC")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	return &Consumer{
		consumer: c,
	}, nil
}

func (c *Consumer) ConsumeEmailJobs(handler func(EmailJob) error) error {
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			// Log error but continue
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		var job EmailJob
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			log.Printf("Failed to unmarshal job: %v\n", err)
			continue
		}

		if err := handler(job); err != nil {
			log.Printf("Failed to handle job: %v\n", err)
			// Continue processing other jobs
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
