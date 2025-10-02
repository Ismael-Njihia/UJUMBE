package utils

import (
	"strings"
)

// ReplaceTemplateVariables replaces template variables like {{variable}} with actual values
func ReplaceTemplateVariables(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, toString(value))
	}
	return result
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
