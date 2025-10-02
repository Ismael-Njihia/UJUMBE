package ses

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESClient struct {
	client *ses.SES
}

func NewSESClient() (*SESClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &SESClient{
		client: ses.New(sess),
	}, nil
}

func (s *SESClient) SendEmail(from, to, subject, htmlBody, textBody string) (string, error) {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
			Body: &ses.Body{},
		},
		Source: aws.String(from),
	}

	if htmlBody != "" {
		input.Message.Body.Html = &ses.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(htmlBody),
		}
	}

	if textBody != "" {
		input.Message.Body.Text = &ses.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(textBody),
		}
	}

	result, err := s.client.SendEmail(input)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return *result.MessageId, nil
}

func (s *SESClient) VerifyDomain(domain string) error {
	_, err := s.client.VerifyDomainIdentity(&ses.VerifyDomainIdentityInput{
		Domain: aws.String(domain),
	})
	if err != nil {
		return fmt.Errorf("failed to verify domain: %w", err)
	}
	return nil
}

func (s *SESClient) VerifyEmail(email string) error {
	_, err := s.client.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	})
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}
