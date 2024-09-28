package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ValidateURL(inputURL string) (string, error) {
	inputURL = strings.TrimSpace(inputURL)

	if inputURL == "" {
		return "", fmt.Errorf("URL is empty")
	}

	// This max length seems to vary depending on the browser and server, but 2048 is common
	if len(inputURL) > 2048 {
		return "", fmt.Errorf("URL is too long (max 2048 characters)")
	}

	if strings.HasSuffix(inputURL, ".") {
		return "", fmt.Errorf("URL cannot end with a dot")
	}

	// Ensure the URL has a scheme
	if !strings.Contains(inputURL, "://") {
		inputURL = "http://" + inputURL
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	switch parsedURL.Scheme {
	case "http", "https", "ftp":
	default:
		return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	parts := strings.Split(parsedURL.Hostname(), ".")

	// If we failed to parse the long URL into at least two parts, something has seriously gone wrong
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid hostname: malformed URL")
	}

	// Subdomain
	if len(parts) == 3 && len(parts[0]) > 63 {
		return "", fmt.Errorf("invalid hostname: subdomain cannot have more than 63 characters")
	}

	// TLD
	tld := parts[len(parts)-1]
	if len(tld) < 2 || len(tld) > 63 {
		return "", fmt.Errorf("invalid hostname: TLD must be between 2 to 63 characters")
	}

	// Hostname
	if len(parsedURL.Hostname()) > 253 {
		return "", fmt.Errorf("URL host is too long (max 253 characters)")
	}

	// Ports
	if portStr := parsedURL.Port(); portStr != "" {
		port, err := strconv.Atoi(portStr) // ASCII to integer
		if err != nil || port < 1 || port > 65535 {
			return "", fmt.Errorf("invalid port number: %s", portStr)
		}
	}

	return parsedURL.String(), nil
}

func GenerateShortID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
