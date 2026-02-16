package utils

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "Test1234",
			wantErr:  false,
		},
		{
			name:     "Long password",
			password: "ThisIsAVeryLongPasswordThatShouldStillWork1234567890",
			wantErr:  false,
		},
		{
			name:     "Short password",
			password: "Abc123",
			wantErr:  false,
		},
		{
			name:     "Password with special characters",
			password: "P@ssw0rd!#$%",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				if hash == tt.password {
					t.Error("HashPassword() returned unhashed password")
				}

				if !strings.HasPrefix(hash, "$2a$") {
					t.Error("HashPassword() did not return bcrypt hash")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		setupPassword string
		expected      bool
	}{
		{
			name:          "Correct password",
			password:      "Test1234",
			setupPassword: "Test1234",
			expected:      true,
		},
		{
			name:          "Incorrect password",
			password:      "WrongPassword",
			setupPassword: "Test1234",
			expected:      false,
		},
		{
			name:          "Empty password against hash",
			password:      "",
			setupPassword: "Test1234",
			expected:      false,
		},
		{
			name:          "Case sensitive",
			password:      "test1234",
			setupPassword: "Test1234",
			expected:      false,
		},
		{
			name:          "Password with spaces",
			password:      "Test 1234",
			setupPassword: "Test 1234",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.setupPassword)
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			result := VerifyPassword(tt.password, hash)

			if result != tt.expected {
				t.Errorf("VerifyPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVerifyPasswordWithInvalidHash(t *testing.T) {
	result := VerifyPassword("Test1234", "invalid-hash")

	if result {
		t.Error("VerifyPassword() should return false for invalid hash")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "Test1234"

	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("First hash failed: %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("Second hash failed: %v", err2)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() should generate different hashes for same password (salt)")
	}

	if !VerifyPassword(password, hash1) {
		t.Error("First hash verification failed")
	}

	if !VerifyPassword(password, hash2) {
		t.Error("Second hash verification failed")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "Test1234"

	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "Test1234"
	hash, _ := HashPassword(password)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hash)
	}
}