# Go Client Example

This example shows how to integrate UJUMBE with your Go application.

## Setup

1. Update the API URL and API key in `main.go`:
   ```go
   const (
       UJUMBE_API_URL = "http://localhost:8080/api/v1"
       API_KEY = "your-api-key-here"
   )
   ```

2. Run the example:
   ```bash
   go run main.go
   ```

## Usage in Your Application

### 1. Copy the UjumbeClient

Copy the `UjumbeClient` struct and methods to your project:

```go
type UjumbeClient struct {
    apiURL string
    apiKey string
    client *http.Client
}
```

### 2. Initialize the Client

```go
client := NewUjumbeClient(
    "http://localhost:8080/api/v1",
    "your-api-key",
)
```

### 3. Send Emails

```go
resp, err := client.SendEmail(EmailRequest{
    From:     "hello@yourdomain.com",
    To:       "user@example.com",
    Subject:  "Welcome!",
    HTMLBody: "<h1>Welcome!</h1>",
})

if err != nil {
    log.Printf("Failed to send email: %v", err)
    return
}

log.Printf("Email sent! ID: %s", resp.EmailID)
```

## Examples Included

1. **Simple Email** - Basic email with subject and body
2. **Template Email** - Email using a pre-defined template
3. **Transactional Email** - Order confirmation example

## Integration Patterns

### User Registration Email

```go
func SendWelcomeEmail(userEmail, userName string) error {
    client := NewUjumbeClient(apiURL, apiKey)
    
    _, err := client.SendEmail(EmailRequest{
        From:    "welcome@yourapp.com",
        To:      userEmail,
        Subject: "Welcome to Our App!",
        HTMLBody: fmt.Sprintf(
            "<h1>Welcome %s!</h1><p>Thanks for joining.</p>",
            userName,
        ),
    })
    
    return err
}
```

### Password Reset Email

```go
func SendPasswordResetEmail(userEmail, resetToken string) error {
    client := NewUjumbeClient(apiURL, apiKey)
    
    resetURL := fmt.Sprintf("https://yourapp.com/reset?token=%s", resetToken)
    
    _, err := client.SendEmail(EmailRequest{
        From:    "noreply@yourapp.com",
        To:      userEmail,
        Subject: "Password Reset Request",
        HTMLBody: fmt.Sprintf(
            "<p>Click <a href='%s'>here</a> to reset your password.</p>"+
            "<p>This link expires in 1 hour.</p>",
            resetURL,
        ),
    })
    
    return err
}
```

### Order Confirmation

```go
type Order struct {
    ID       string
    Items    []string
    Total    float64
    Customer string
}

func SendOrderConfirmation(order Order, customerEmail string) error {
    client := NewUjumbeClient(apiURL, apiKey)
    
    itemsList := ""
    for _, item := range order.Items {
        itemsList += fmt.Sprintf("<li>%s</li>", item)
    }
    
    htmlBody := fmt.Sprintf(`
        <h1>Order Confirmed!</h1>
        <p>Order ID: %s</p>
        <h2>Items:</h2>
        <ul>%s</ul>
        <p>Total: $%.2f</p>
    `, order.ID, itemsList, order.Total)
    
    _, err := client.SendEmail(EmailRequest{
        From:     "orders@yourapp.com",
        To:       customerEmail,
        Subject:  fmt.Sprintf("Order Confirmation #%s", order.ID),
        HTMLBody: htmlBody,
    })
    
    return err
}
```

## Error Handling

Always handle errors appropriately:

```go
resp, err := client.SendEmail(req)
if err != nil {
    // Log the error
    log.Printf("Failed to send email: %v", err)
    
    // Optionally retry with exponential backoff
    // Or queue for later retry
    
    return err
}

// Email sent successfully
log.Printf("Email sent: %s", resp.EmailID)
```

## Best Practices

1. **Initialize once** - Create the client once and reuse it
2. **Handle errors** - Always check for errors
3. **Use templates** - For consistent branding
4. **Verify domains** - Before sending production emails
5. **Monitor quota** - Check `resp.Remaining` to track usage
6. **Async sending** - Consider sending emails in background goroutines

## Advanced Usage

### Concurrent Email Sending

```go
func SendBulkEmails(emails []string) {
    client := NewUjumbeClient(apiURL, apiKey)
    
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 10) // Limit to 10 concurrent sends
    
    for _, email := range emails {
        wg.Add(1)
        go func(to string) {
            defer wg.Done()
            semaphore <- struct{}{}        // Acquire
            defer func() { <-semaphore }() // Release
            
            _, err := client.SendEmail(EmailRequest{
                From:     "newsletter@yourapp.com",
                To:       to,
                Subject:  "Monthly Newsletter",
                HTMLBody: "<h1>Newsletter</h1>",
            })
            
            if err != nil {
                log.Printf("Failed to send to %s: %v", to, err)
            }
        }(email)
    }
    
    wg.Wait()
}
```

### With Retry Logic

```go
func SendEmailWithRetry(req EmailRequest, maxRetries int) error {
    client := NewUjumbeClient(apiURL, apiKey)
    
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        _, err := client.SendEmail(req)
        if err == nil {
            return nil
        }
        
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
    }
    
    return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

## Testing

For testing, you can mock the client:

```go
type EmailSender interface {
    SendEmail(EmailRequest) (*EmailResponse, error)
}

type MockEmailSender struct{}

func (m *MockEmailSender) SendEmail(req EmailRequest) (*EmailResponse, error) {
    return &EmailResponse{
        Success:   true,
        EmailID:   "test-id",
        Message:   "Test email",
        Remaining: 99,
    }, nil
}
```
