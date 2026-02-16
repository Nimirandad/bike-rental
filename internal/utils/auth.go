package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func ValidateAdminBasicAuth(authHeader string) error {
	if authHeader == "" {
		return fmt.Errorf("authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return fmt.Errorf("authorization header must use Basic authentication")
	}

	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")

	expectedCredentials := os.Getenv("ADMIN_CREDENTIALS")
	if expectedCredentials == "" {
		return fmt.Errorf("admin credentials not configured")
	}

	if encodedCredentials != expectedCredentials {
		return fmt.Errorf("invalid admin credentials")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return fmt.Errorf("invalid authorization format")
	}

	credentials := string(decodedBytes)
	if !strings.Contains(credentials, ":") {
		return fmt.Errorf("invalid credentials format")
	}

	return nil
}