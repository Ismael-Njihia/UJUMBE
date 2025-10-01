"""
UJUMBE Python Client Example

A simple Python client for sending emails via UJUMBE API.
"""

import requests
from typing import Optional, Dict, Any
import json


class UjumbeClient:
    """Client for interacting with UJUMBE email API."""
    
    def __init__(self, api_url: str, api_key: str):
        """
        Initialize UJUMBE client.
        
        Args:
            api_url: Base URL for UJUMBE API (e.g., http://localhost:8080/api/v1)
            api_key: Your UJUMBE API key
        """
        self.api_url = api_url.rstrip('/')
        self.api_key = api_key
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'X-API-Key': api_key
        })
    
    def send_email(
        self,
        from_email: str,
        to_email: str,
        subject: Optional[str] = None,
        html_body: Optional[str] = None,
        text_body: Optional[str] = None,
        template_id: Optional[str] = None,
        template_data: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """
        Send an email via UJUMBE.
        
        Args:
            from_email: Sender email address (must be from verified domain)
            to_email: Recipient email address
            subject: Email subject (required if not using template)
            html_body: HTML email body
            text_body: Plain text email body
            template_id: Template UUID to use
            template_data: Data to fill template variables
            
        Returns:
            Response dict with success, email_id, message, and remaining quota
            
        Raises:
            requests.HTTPError: If the API request fails
        """
        payload = {
            'from': from_email,
            'to': to_email
        }
        
        if template_id:
            payload['template_id'] = template_id
            if template_data:
                payload['template_data'] = template_data
        else:
            if not subject:
                raise ValueError("subject is required when not using a template")
            payload['subject'] = subject
            if html_body:
                payload['html_body'] = html_body
            if text_body:
                payload['text_body'] = text_body
        
        response = self.session.post(
            f'{self.api_url}/emails/send',
            json=payload
        )
        response.raise_for_status()
        return response.json()
    
    def get_email_status(self, email_id: str) -> Dict[str, Any]:
        """Get the status of a sent email."""
        response = self.session.get(f'{self.api_url}/emails/{email_id}')
        response.raise_for_status()
        return response.json()
    
    def get_email_logs(self, email_id: str) -> list:
        """Get logs for a sent email."""
        response = self.session.get(f'{self.api_url}/emails/{email_id}/logs')
        response.raise_for_status()
        return response.json()
    
    def get_analytics(self) -> Dict[str, Any]:
        """Get email analytics."""
        response = self.session.get(f'{self.api_url}/analytics')
        response.raise_for_status()
        return response.json()
    
    def get_quota(self) -> Dict[str, Any]:
        """Get remaining email quota."""
        response = self.session.get(f'{self.api_url}/quota')
        response.raise_for_status()
        return response.json()
    
    def create_template(
        self,
        name: str,
        subject: str,
        html_body: str,
        text_body: Optional[str] = None
    ) -> Dict[str, Any]:
        """Create a new email template."""
        payload = {
            'name': name,
            'subject': subject,
            'html_body': html_body
        }
        if text_body:
            payload['text_body'] = text_body
        
        response = self.session.post(
            f'{self.api_url}/templates',
            json=payload
        )
        response.raise_for_status()
        return response.json()
    
    def list_templates(self) -> list:
        """List all templates."""
        response = self.session.get(f'{self.api_url}/templates')
        response.raise_for_status()
        return response.json()
    
    def add_domain(self, domain: str) -> Dict[str, Any]:
        """Add a domain for verification."""
        response = self.session.post(
            f'{self.api_url}/domains',
            json={'domain': domain}
        )
        response.raise_for_status()
        return response.json()
    
    def list_domains(self) -> list:
        """List all domains."""
        response = self.session.get(f'{self.api_url}/domains')
        response.raise_for_status()
        return response.json()


def main():
    # Initialize client
    # Update these values with your actual API URL and key
    API_URL = 'http://localhost:8080/api/v1'
    API_KEY = 'your-api-key-here'
    
    client = UjumbeClient(API_URL, API_KEY)
    
    # Example 1: Send a simple email
    print("=== Example 1: Simple Email ===")
    try:
        response = client.send_email(
            from_email='hello@yourdomain.com',
            to_email='user@example.com',
            subject='Welcome to Our Service',
            html_body='<h1>Welcome!</h1><p>Thank you for signing up.</p>',
            text_body='Welcome! Thank you for signing up.'
        )
        print(f"✓ Email sent! ID: {response['email_id']}, Remaining: {response['remaining']}")
    except requests.HTTPError as e:
        print(f"✗ Error: {e.response.json()}")
    
    # Example 2: Send email with template
    print("\n=== Example 2: Template Email ===")
    try:
        response = client.send_email(
            from_email='hello@yourdomain.com',
            to_email='user@example.com',
            template_id='your-template-uuid',  # Replace with actual template ID
            template_data={
                'name': 'John Doe',
                'company': 'ACME Inc',
                'action_url': 'https://example.com/verify'
            }
        )
        print(f"✓ Email sent! ID: {response['email_id']}, Remaining: {response['remaining']}")
    except requests.HTTPError as e:
        print(f"✗ Error: {e.response.json()}")
    except ValueError as e:
        print(f"✗ Note: {e}")
    
    # Example 3: Check quota
    print("\n=== Example 3: Check Quota ===")
    try:
        quota = client.get_quota()
        print(f"Free emails remaining: {quota['free_emails_remaining']}")
        print(f"Paid email balance: {quota['paid_emails_balance']}")
    except requests.HTTPError as e:
        print(f"✗ Error: {e.response.json()}")
    
    # Example 4: Get analytics
    print("\n=== Example 4: Analytics ===")
    try:
        analytics = client.get_analytics()
        print(f"Total sent: {analytics['total_emails_sent']}")
        print(f"Total failed: {analytics['total_emails_failed']}")
        print(f"Success rate: {analytics['success_rate']}%")
    except requests.HTTPError as e:
        print(f"✗ Error: {e.response.json()}")


if __name__ == '__main__':
    main()
