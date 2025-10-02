# Python Client Example

This example shows how to integrate UJUMBE with your Python application.

## Requirements

```bash
pip install requests
```

## Setup

Update the API URL and API key in `ujumbe_client.py`:

```python
API_URL = 'http://localhost:8080/api/v1'
API_KEY = 'your-api-key-here'
```

## Run the Example

```bash
python ujumbe_client.py
```

## Usage in Your Application

### 1. Copy the UjumbeClient class

Copy the `UjumbeClient` class to your project.

### 2. Initialize the Client

```python
from ujumbe_client import UjumbeClient

client = UjumbeClient(
    api_url='http://localhost:8080/api/v1',
    api_key='your-api-key'
)
```

### 3. Send Emails

```python
response = client.send_email(
    from_email='hello@yourdomain.com',
    to_email='user@example.com',
    subject='Welcome!',
    html_body='<h1>Welcome!</h1>'
)

print(f"Email sent! ID: {response['email_id']}")
```

## Integration Examples

### Django Integration

```python
# emails.py
from django.conf import settings
from ujumbe_client import UjumbeClient

def get_ujumbe_client():
    return UjumbeClient(
        api_url=settings.UJUMBE_API_URL,
        api_key=settings.UJUMBE_API_KEY
    )

def send_welcome_email(user):
    client = get_ujumbe_client()
    
    try:
        response = client.send_email(
            from_email='welcome@yourapp.com',
            to_email=user.email,
            subject='Welcome to Our App!',
            html_body=f'<h1>Welcome {user.first_name}!</h1><p>Thanks for joining.</p>'
        )
        return response['email_id']
    except Exception as e:
        # Log error
        print(f"Failed to send welcome email: {e}")
        return None
```

### Flask Integration

```python
# app.py
from flask import Flask
from ujumbe_client import UjumbeClient

app = Flask(__name__)
ujumbe = UjumbeClient(
    api_url=app.config['UJUMBE_API_URL'],
    api_key=app.config['UJUMBE_API_KEY']
)

@app.route('/signup', methods=['POST'])
def signup():
    # ... user registration logic ...
    
    # Send welcome email
    try:
        ujumbe.send_email(
            from_email='noreply@yourapp.com',
            to_email=user.email,
            subject='Welcome!',
            html_body='<h1>Welcome to our app!</h1>'
        )
    except Exception as e:
        app.logger.error(f"Failed to send email: {e}")
    
    return {'success': True}
```

### Celery Task (Async)

```python
# tasks.py
from celery import Celery
from ujumbe_client import UjumbeClient

app = Celery('tasks')
client = UjumbeClient(api_url='...', api_key='...')

@app.task
def send_email_async(from_email, to_email, subject, html_body):
    """Send email asynchronously."""
    try:
        response = client.send_email(
            from_email=from_email,
            to_email=to_email,
            subject=subject,
            html_body=html_body
        )
        return response['email_id']
    except Exception as e:
        # Retry on failure
        raise self.retry(exc=e, countdown=60)

# Usage
send_email_async.delay(
    from_email='noreply@yourapp.com',
    to_email='user@example.com',
    subject='Hello',
    html_body='<h1>Hello!</h1>'
)
```

## Common Use Cases

### Password Reset Email

```python
def send_password_reset_email(user_email, reset_token):
    client = UjumbeClient(api_url='...', api_key='...')
    
    reset_url = f'https://yourapp.com/reset?token={reset_token}'
    
    html_body = f"""
    <h1>Password Reset Request</h1>
    <p>Click the link below to reset your password:</p>
    <p><a href="{reset_url}">Reset Password</a></p>
    <p>This link expires in 1 hour.</p>
    """
    
    return client.send_email(
        from_email='noreply@yourapp.com',
        to_email=user_email,
        subject='Password Reset Request',
        html_body=html_body
    )
```

### Order Confirmation

```python
def send_order_confirmation(order):
    client = UjumbeClient(api_url='...', api_key='...')
    
    items_html = ''.join([
        f'<li>{item.name} - ${item.price}</li>'
        for item in order.items
    ])
    
    html_body = f"""
    <h1>Order Confirmed!</h1>
    <p>Order ID: {order.id}</p>
    <h2>Items:</h2>
    <ul>{items_html}</ul>
    <p>Total: ${order.total}</p>
    <p><a href="https://yourapp.com/orders/{order.id}">View Order</a></p>
    """
    
    return client.send_email(
        from_email='orders@yourapp.com',
        to_email=order.customer_email,
        subject=f'Order Confirmation #{order.id}',
        html_body=html_body
    )
```

### Newsletter with Template

```python
def send_newsletter(subscribers):
    client = UjumbeClient(api_url='...', api_key='...')
    
    # Get template ID (create template first via dashboard or API)
    template_id = 'your-newsletter-template-id'
    
    for subscriber in subscribers:
        try:
            client.send_email(
                from_email='newsletter@yourapp.com',
                to_email=subscriber.email,
                template_id=template_id,
                template_data={
                    'name': subscriber.name,
                    'unsubscribe_url': f'https://yourapp.com/unsubscribe/{subscriber.token}'
                }
            )
        except Exception as e:
            print(f"Failed to send to {subscriber.email}: {e}")
```

## Error Handling

```python
import requests
from requests.exceptions import HTTPError, RequestException

try:
    response = client.send_email(
        from_email='hello@yourdomain.com',
        to_email='user@example.com',
        subject='Test',
        html_body='<h1>Test</h1>'
    )
except HTTPError as e:
    # Handle HTTP errors (4xx, 5xx)
    error_data = e.response.json()
    print(f"API Error: {error_data.get('error')}")
except RequestException as e:
    # Handle connection errors
    print(f"Connection Error: {e}")
except ValueError as e:
    # Handle validation errors
    print(f"Validation Error: {e}")
```

## Bulk Sending with Concurrency

```python
from concurrent.futures import ThreadPoolExecutor, as_completed

def send_bulk_emails(recipients, subject, html_body):
    client = UjumbeClient(api_url='...', api_key='...')
    
    def send_one(to_email):
        try:
            response = client.send_email(
                from_email='newsletter@yourapp.com',
                to_email=to_email,
                subject=subject,
                html_body=html_body
            )
            return (to_email, response['email_id'], None)
        except Exception as e:
            return (to_email, None, str(e))
    
    results = []
    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = [executor.submit(send_one, email) for email in recipients]
        
        for future in as_completed(futures):
            results.append(future.result())
    
    # Process results
    success = [r for r in results if r[2] is None]
    failed = [r for r in results if r[2] is not None]
    
    print(f"Sent: {len(success)}, Failed: {len(failed)}")
    return results
```

## Testing

Mock the client for testing:

```python
from unittest.mock import Mock

def test_send_welcome_email():
    # Mock the client
    mock_client = Mock(spec=UjumbeClient)
    mock_client.send_email.return_value = {
        'success': True,
        'email_id': 'test-id',
        'message': 'Email queued',
        'remaining': 99
    }
    
    # Test your function
    response = mock_client.send_email(
        from_email='test@example.com',
        to_email='user@example.com',
        subject='Test',
        html_body='<h1>Test</h1>'
    )
    
    assert response['success'] == True
    assert response['email_id'] == 'test-id'
```

## Best Practices

1. **Use environment variables** for API key and URL
2. **Handle errors gracefully** - don't let email failures break your app
3. **Send emails asynchronously** - use Celery or similar
4. **Monitor quota** - check `response['remaining']`
5. **Use templates** - for consistent branding
6. **Verify domains** - before production use
7. **Log failures** - for debugging and monitoring
