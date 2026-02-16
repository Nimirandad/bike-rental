package utils

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestValidateAdminBasicAuth(t *testing.T) {
	originalCreds := os.Getenv("ADMIN_CREDENTIALS")
	defer func() {
		if originalCreds != "" {
			os.Setenv("ADMIN_CREDENTIALS", originalCreds)
		} else {
			os.Unsetenv("ADMIN_CREDENTIALS")
		}
	}()

	t.Run("Valid credentials", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		authHeader := "Basic " + encodedCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err != nil {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want nil", err)
		}
	})

	t.Run("Valid credentials with extra spaces", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		authHeader := "Basic " + encodedCreds + "  "

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for extra spaces, got nil")
		}
	})

	t.Run("Missing authorization header", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		err := ValidateAdminBasicAuth("")
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for empty header, got nil")
		}

		expectedError := "authorization header is required"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Invalid header format - no Basic prefix", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		authHeader := encodedCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for missing Basic prefix, got nil")
		}

		expectedError := "authorization header must use Basic authentication"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Invalid header format - wrong case", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		authHeader := "basic " + encodedCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for wrong case, got nil")
		}
	})

	t.Run("Invalid base64 encoding", func(t *testing.T) {
		encodedCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", encodedCreds)

		authHeader := "Basic not-valid-base64!"

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for invalid base64, got nil")
		}

		expectedError := "invalid admin credentials"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Invalid credentials format - no colon", func(t *testing.T) {
		wrongCredentials := base64.StdEncoding.EncodeToString([]byte("adminpassword123"))
		os.Setenv("ADMIN_CREDENTIALS", wrongCredentials)

		authHeader := "Basic " + wrongCredentials

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for missing colon, got nil")
		}

		expectedError := "invalid credentials format"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Wrong username", func(t *testing.T) {
		correctCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", correctCreds)

		wrongCreds := base64.StdEncoding.EncodeToString([]byte("wronguser:password123"))
		authHeader := "Basic " + wrongCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for wrong username, got nil")
		}

		expectedError := "invalid admin credentials"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		correctCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", correctCreds)

		wrongCreds := base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword"))
		authHeader := "Basic " + wrongCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for wrong password, got nil")
		}

		expectedError := "invalid admin credentials"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Missing ADMIN_CREDENTIALS env var", func(t *testing.T) {
		os.Unsetenv("ADMIN_CREDENTIALS")

		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password123"))

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for missing env var, got nil")
		}

		expectedError := "admin credentials not configured"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("ADMIN_CREDENTIALS without colon", func(t *testing.T) {
		wrongCreds := base64.StdEncoding.EncodeToString([]byte("adminpassword"))
		os.Setenv("ADMIN_CREDENTIALS", wrongCreds)

		authHeader := "Basic " + wrongCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for malformed env var, got nil")
		}

		expectedError := "invalid credentials format"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Empty credentials", func(t *testing.T) {
		correctCreds := base64.StdEncoding.EncodeToString([]byte("admin:password123"))
		os.Setenv("ADMIN_CREDENTIALS", correctCreds)

		emptyCreds := base64.StdEncoding.EncodeToString([]byte(":"))
		authHeader := "Basic " + emptyCreds

		err := ValidateAdminBasicAuth(authHeader)
		if err == nil {
			t.Error("ValidateAdminBasicAuth() expected error for empty credentials, got nil")
		}

		expectedError := "invalid admin credentials"
		if err.Error() != expectedError {
			t.Errorf("ValidateAdminBasicAuth() error = %v, want %v", err.Error(), expectedError)
		}
	})
}