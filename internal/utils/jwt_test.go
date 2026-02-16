package utils

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	t.Run("Valid user", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		token, err := GenerateJWT(user)
		if err != nil {
			t.Errorf("GenerateJWT() error = %v, want nil", err)
		}

		if token == "" {
			t.Error("GenerateJWT() returned empty token")
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			t.Errorf("GenerateJWT() token has %d parts, want 3", len(parts))
		}
	})

	t.Run("Missing JWT_SECRET", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")

		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		token, err := GenerateJWT(user)
		if err == nil {
			t.Error("GenerateJWT() expected error, got nil")
		}

		if token != "" {
			t.Error("GenerateJWT() expected empty token on error")
		}

		expectedError := "JWT_SECRET environment variable not set"
		if err.Error() != expectedError {
			t.Errorf("GenerateJWT() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Token contains correct claims", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		user := &models.User{
			ID:        42,
			Email:     "jane@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
		}

		tokenString, err := GenerateJWT(user)
		if err != nil {
			t.Fatalf("GenerateJWT() error = %v", err)
		}

		token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key"), nil
		})

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			t.Fatal("Failed to parse claims")
		}

		if claims.Sub != user.ID {
			t.Errorf("Claims.Sub = %v, want %v", claims.Sub, user.ID)
		}
		if claims.Email != user.Email {
			t.Errorf("Claims.Email = %v, want %v", claims.Email, user.Email)
		}
		if claims.FirstName != user.FirstName {
			t.Errorf("Claims.FirstName = %v, want %v", claims.FirstName, user.FirstName)
		}
		if claims.LastName != user.LastName {
			t.Errorf("Claims.LastName = %v, want %v", claims.LastName, user.LastName)
		}
	})
}

func TestValidateJWT(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	t.Run("Valid token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		tokenString, err := GenerateJWT(user)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			t.Errorf("ValidateJWT() error = %v, want nil", err)
		}

		if claims == nil {
			t.Fatal("ValidateJWT() claims = nil, want non-nil")
		}

		if claims.Sub != user.ID {
			t.Errorf("Claims.Sub = %v, want %v", claims.Sub, user.ID)
		}
		if claims.Email != user.Email {
			t.Errorf("Claims.Email = %v, want %v", claims.Email, user.Email)
		}
	})

	t.Run("Invalid token - malformed", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		claims, err := ValidateJWT("invalid.token.string")
		if err == nil {
			t.Error("ValidateJWT() expected error for malformed token, got nil")
		}

		if claims != nil {
			t.Error("ValidateJWT() expected nil claims for malformed token")
		}
	})

	t.Run("Invalid token - wrong signature", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		tokenString, _ := GenerateJWT(user)

		os.Setenv("JWT_SECRET", "different-secret-key")

		claims, err := ValidateJWT(tokenString)
		if err == nil {
			t.Error("ValidateJWT() expected error for wrong signature, got nil")
		}

		if claims != nil {
			t.Error("ValidateJWT() expected nil claims for wrong signature")
		}
	})

	t.Run("Missing JWT_SECRET", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")

		claims, err := ValidateJWT("some.token.string")
		if err == nil {
			t.Error("ValidateJWT() expected error, got nil")
		}

		if claims != nil {
			t.Error("ValidateJWT() expected nil claims")
		}

		expectedError := "JWT_SECRET environment variable not set"
		if err.Error() != expectedError {
			t.Errorf("ValidateJWT() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		expirationTime := time.Now().Add(-1 * time.Hour)

		claims := JWTClaims{
			Sub:       1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		parsedClaims, err := ValidateJWT(tokenString)
		if err == nil {
			t.Error("ValidateJWT() expected error for expired token, got nil")
		}

		if parsedClaims != nil {
			t.Error("ValidateJWT() expected nil claims for expired token")
		}
	})
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantError   bool
		errorString string
	}{
		{
			name:       "Valid Bearer token",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantError:  false,
		},
		{
			name:       "Valid Bearer token with extra spaces",
			authHeader: "Bearer   token-with-spaces   ",
			wantToken:  "token-with-spaces",
			wantError:  false,
		},
		{
			name:       "Token with double quotes",
			authHeader: "Bearer \"token-in-quotes\"",
			wantToken:  "token-in-quotes",
			wantError:  false,
		},
		{
			name:       "Token with single quotes",
			authHeader: "Bearer 'token-in-quotes'",
			wantToken:  "token-in-quotes",
			wantError:  false,
		},
		{
			name:        "Empty header",
			authHeader:  "",
			wantToken:   "",
			wantError:   true,
			errorString: "authorization header is required",
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  "token-without-bearer",
			wantToken:   "",
			wantError:   true,
			errorString: "authorization header format must be Bearer {token}",
		},
		{
			name:        "Wrong case Bearer",
			authHeader:  "bearer token",
			wantToken:   "",
			wantError:   true,
			errorString: "authorization header format must be Bearer {token}",
		},
		{
			name:        "Only Bearer without token",
			authHeader:  "Bearer",
			wantToken:   "",
			wantError:   true,
			errorString: "authorization header format must be Bearer {token}",
		},
		{
			name:        "Bearer with only space",
			authHeader:  "Bearer ",
			wantToken:   "",
			wantError:   false,
			errorString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tt.authHeader)

			if tt.wantError && err == nil {
				t.Errorf("ExtractTokenFromHeader() expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("ExtractTokenFromHeader() unexpected error = %v", err)
			}

			if tt.wantError && err != nil && tt.errorString != "" {
				if err.Error() != tt.errorString {
					t.Errorf("ExtractTokenFromHeader() error = %v, want %v", err.Error(), tt.errorString)
				}
			}

			if token != tt.wantToken {
				t.Errorf("ExtractTokenFromHeader() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}

func BenchmarkGenerateJWT(b *testing.B) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateJWT(user)
	}
}

func BenchmarkValidateJWT(b *testing.B) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	tokenString, _ := GenerateJWT(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateJWT(tokenString)
	}
}