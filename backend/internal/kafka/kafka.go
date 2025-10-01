package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type EmailJob struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	TemplateID string            `json:"template_id,omitempty"`
	From       string            `json:"from"`
	To         string            `json:"to"`
	Subject    string            `json:"subject"`
	HTMLBody   string            `json:"html_body"`
	TextBody   string            `json:"text_body"`
	Variables  map[string]string `json:"variables,omitempty"`
}

type Producer struct {
	writer *kafka.Writer
}

type Consumer struct {
	reader *kafka.Reader
}

func NewProducer(brokers, topic string) *Producer {
	brokerList := strings.Split(brokers, ",")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerList...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    10,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{writer: writer}
}

func (p *Producer) PublishEmail(ctx context.Context, job *EmailJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal email job: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(job.UserID),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func NewConsumer(brokers, topic, groupID string) *Consumer {
	brokerList := strings.Split(brokers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokerList,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) ConsumeEmails(ctx context.Context, handler func(*EmailJob) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			var job EmailJob
			if err := json.Unmarshal(msg.Value, &job); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				c.reader.CommitMessages(ctx, msg)
				continue
			}

			if err := handler(&job); err != nil {
				log.Printf("Error handling email job: %v", err)
				// Don't commit on error - message will be retried
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
