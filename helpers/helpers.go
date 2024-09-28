package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ValidateURL(inputURL string) (string, error) {
	inputURL = strings.TrimSpace(inputURL)

	if inputURL == "" {
		return "", errors.New("URL is empty")
	}

	// This max length seems to vary depending on the browser and server, but 2048 is common
	if len(inputURL) > 2048 {
		return "", errors.New("URL is too long (max 2048 characters)")
	}

	// Go's URL.Parse() should catch most invalid URLs
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	switch parsedURL.Scheme {
	case "http", "https", "ftp":
	default:
		return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	if len(parsedURL.Host) > 253 {
		return "", errors.New("URL host is too long (max 253 characters)")
	}

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
