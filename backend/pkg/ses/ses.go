package ses

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESClient struct {
	client *ses.SES
}

type EmailMessage struct {
	From     string
	To       []string
	Subject  string
	HTMLBody string
	TextBody string
}

func NewSESClient(region, accessKeyID, secretAccessKey string) (*SESClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &SESClient{
		client: ses.New(sess),
	}, nil
}

func (c *SESClient) SendEmail(msg *EmailMessage) (string, error) {
	input := &ses.SendEmailInput{
		Source: aws.String(msg.From),
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(msg.To),
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data:    aws.String(msg.Subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &ses.Body{},
		},
	}

	if msg.HTMLBody != "" {
		input.Message.Body.Html = &ses.Content{
			Data:    aws.String(msg.HTMLBody),
			Charset: aws.String("UTF-8"),
		}
	}

	if msg.TextBody != "" {
		input.Message.Body.Text = &ses.Content{
			Data:    aws.String(msg.TextBody),
			Charset: aws.String("UTF-8"),
		}
	}

	result, err := c.client.SendEmail(input)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return *result.MessageId, nil
}

func (c *SESClient) VerifyEmailIdentity(email string) error {
	_, err := c.client.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	})
	return err
}

func (c *SESClient) VerifyDomainIdentity(domain string) (string, error) {
	result, err := c.client.VerifyDomainIdentity(&ses.VerifyDomainIdentityInput{
		Domain: aws.String(domain),
	})
	if err != nil {
		return "", err
	}
	return *result.VerificationToken, nil
}

func (c *SESClient) GetSendQuota() (*ses.GetSendQuotaOutput, error) {
	return c.client.GetSendQuota(&ses.GetSendQuotaInput{})
}

func (c *SESClient) GetSendStatistics() (*ses.GetSendStatisticsOutput, error) {
	return c.client.GetSendStatistics(&ses.GetSendStatisticsInput{})
}
