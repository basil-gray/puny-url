package helpers

import (
	"puny-url/internal/helpers"
	"strings"
	"testing"
)

func TestValidURLs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "URL without prefix",
			input:       "example.com",
			expected:    "http://example.com",
			expectError: false,
		},
		{
			name:        "URL with http prefix",
			input:       "http://example.com",
			expected:    "http://example.com",
			expectError: false,
		},
		{
			name:        "URL with https prefix",
			input:       "https://example.com",
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "URL with path",
			input:       "example.com/path",
			expected:    "http://example.com/path",
			expectError: false,
		},
		{
			name:        "URL with query parameters",
			input:       "example.com?param=value",
			expected:    "http://example.com?param=value",
			expectError: false,
		},
		{
			name:        "URL with percent encoding",
			input:       "http://example.com/path%20with%20spaces",
			expected:    "http://example.com/path%20with%20spaces",
			expectError: false,
		},
		{
			name:        "URL with port",
			input:       "example.com:8080",
			expected:    "http://example.com:8080",
			expectError: false,
		},
		{
			name:        "TLD with 2 characters",
			input:       "http://example.co",
			expected:    "http://example.co",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.ParseURL(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateURL(%q) expected an error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL(%q) returned an unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ValidateURL(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestInvalidURLs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "URL with unsupported scheme",
			input:       "asdf://example.com:8080",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Malformed URL with double prefix",
			input:       "http://http://example.com:8080",
			expected:    "",
			expectError: true,
		},
		{
			name:        "TLD more than 63 characters",
			input:       "http://example.com" + strings.Repeat("m", 63),
			expected:    "",
			expectError: true,
		},
		{
			name:        "Subdomain more than 63 characters",
			input:       "http://www" + strings.Repeat("w", 63) + ".example.com",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid port",
			input:       "http://www.example.com:100000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "URL that ends with a .",
			input:       "http://www.example.",
			expected:    "",
			expectError: true,
		},
		{
			name:        "URL that's more than 2048 characters",
			input:       "http://www.example.com/" + strings.Repeat("a", 2048),
			expected:    "",
			expectError: true,
		},
		{
			name:        "URL with spaces",
			input:       "http://exam ple.com",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.ParseURL(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateURL(%q) expected an error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL(%q) returned an unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ValidateURL(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}
