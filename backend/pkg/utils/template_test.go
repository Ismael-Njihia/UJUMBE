package utils

import (
	"testing"
)

func TestReplaceTemplateVariables(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "Single variable",
			template: "Hello {{name}}!",
			data:     map[string]interface{}{"name": "John"},
			expected: "Hello John!",
		},
		{
			name:     "Multiple variables",
			template: "Hello {{name}}, welcome to {{company}}!",
			data:     map[string]interface{}{"name": "John", "company": "ACME"},
			expected: "Hello John, welcome to ACME!",
		},
		{
			name:     "No variables",
			template: "Hello World!",
			data:     map[string]interface{}{},
			expected: "Hello World!",
		},
		{
			name:     "Missing variable data",
			template: "Hello {{name}}!",
			data:     map[string]interface{}{"other": "value"},
			expected: "Hello {{name}}!",
		},
		{
			name:     "Empty template",
			template: "",
			data:     map[string]interface{}{"name": "John"},
			expected: "",
		},
		{
			name:     "Variable in HTML",
			template: "<h1>Welcome {{name}}!</h1><p>Email: {{email}}</p>",
			data:     map[string]interface{}{"name": "John", "email": "john@example.com"},
			expected: "<h1>Welcome John!</h1><p>Email: john@example.com</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceTemplateVariables(tt.template, tt.data)
			if result != tt.expected {
				t.Errorf("ReplaceTemplateVariables() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReplaceTemplateVariablesWithNilData(t *testing.T) {
	template := "Hello {{name}}!"
	result := ReplaceTemplateVariables(template, nil)
	
	if result != template {
		t.Errorf("Expected template to remain unchanged with nil data, got %v", result)
	}
}

func TestReplaceTemplateVariablesWithNonStringValues(t *testing.T) {
	template := "Count: {{count}}, Active: {{active}}"
	data := map[string]interface{}{
		"count":  123,
		"active": true,
	}
	
	result := ReplaceTemplateVariables(template, data)
	
	// Non-string values should be ignored or converted
	if result == "" {
		t.Error("Result should not be empty")
	}
}
