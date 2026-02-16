package utils

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantValid   bool
		wantMessage string
	}{
		{
			name:        "Valid email",
			email:       "test@example.com",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid email with subdomain",
			email:       "user@mail.example.com",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid email with plus",
			email:       "user+tag@example.com",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid email with numbers",
			email:       "user123@example456.com",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Empty email",
			email:       "",
			wantValid:   false,
			wantMessage: "Email is required",
		},
		{
			name:        "Email with spaces",
			email:       "  test@example.com  ",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Invalid format - no @",
			email:       "testexample.com",
			wantValid:   false,
			wantMessage: "Invalid email format",
		},
		{
			name:        "Invalid format - no domain",
			email:       "test@",
			wantValid:   false,
			wantMessage: "Invalid email format",
		},
		{
			name:        "Invalid format - no TLD",
			email:       "test@example",
			wantValid:   false,
			wantMessage: "Invalid email format",
		},
		{
			name:        "Invalid format - spaces in email",
			email:       "test user@example.com",
			wantValid:   false,
			wantMessage: "Invalid email format",
		},
		{
			name:        "Invalid format - multiple @",
			email:       "test@@example.com",
			wantValid:   false,
			wantMessage: "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, message := ValidateEmail(tt.email)

			if valid != tt.wantValid {
				t.Errorf("ValidateEmail() valid = %v, want %v", valid, tt.wantValid)
			}

			if message != tt.wantMessage {
				t.Errorf("ValidateEmail() message = %v, want %v", message, tt.wantMessage)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantValid   bool
		wantMessage string
	}{
		{
			name:        "Valid password",
			password:    "Password123",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid password - minimum length",
			password:    "Pass123w",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid password with special chars",
			password:    "P@ssw0rd!",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Empty password",
			password:    "",
			wantValid:   false,
			wantMessage: "Password is required",
		},
		{
			name:        "Too short",
			password:    "Pass12",
			wantValid:   false,
			wantMessage: "Password must be at least 8 characters",
		},
		{
			name:        "No numbers",
			password:    "PasswordOnly",
			wantValid:   false,
			wantMessage: "Password must contain at least one letter and one number",
		},
		{
			name:        "No letters",
			password:    "12345678",
			wantValid:   false,
			wantMessage: "Password must contain at least one letter and one number",
		},
		{
			name:        "Too long (>100 chars)",
			password:    "Password123" + string(make([]byte, 100)),
			wantValid:   false,
			wantMessage: "Password is too long (max 100 characters)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, message := ValidatePassword(tt.password)

			if valid != tt.wantValid {
				t.Errorf("ValidatePassword() valid = %v, want %v", valid, tt.wantValid)
			}

			if message != tt.wantMessage {
				t.Errorf("ValidatePassword() message = %v, want %v", message, tt.wantMessage)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		fieldName   string
		wantValid   bool
		wantMessage string
	}{
		{
			name:        "Valid first name",
			inputName:   "John",
			fieldName:   "First name",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid name with hyphen",
			inputName:   "Mary-Jane",
			fieldName:   "First name",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Valid name with spaces",
			inputName:   "John Smith",
			fieldName:   "Last name",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Empty name",
			inputName:   "",
			fieldName:   "First name",
			wantValid:   false,
			wantMessage: "First name is required",
		},
		{
			name:        "Name with spaces trimmed",
			inputName:   "  John  ",
			fieldName:   "First name",
			wantValid:   true,
			wantMessage: "",
		},
		{
			name:        "Too short (1 char)",
			inputName:   "J",
			fieldName:   "First name",
			wantValid:   false,
			wantMessage: "First name must be at least 2 characters",
		},
		{
			name:        "Too long (>50 chars)",
			inputName:   string(make([]byte, 51)),
			fieldName:   "First name",
			wantValid:   false,
			wantMessage: "First name is too long (max 50 characters)",
		},
		{
			name:        "Contains numbers",
			inputName:   "John123",
			fieldName:   "First name",
			wantValid:   false,
			wantMessage: "First name can only contain letters, spaces, and hyphens",
		},
		{
			name:        "Contains special characters",
			inputName:   "John@Smith",
			fieldName:   "First name",
			wantValid:   false,
			wantMessage: "First name can only contain letters, spaces, and hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, message := ValidateName(tt.inputName, tt.fieldName)

			if valid != tt.wantValid {
				t.Errorf("ValidateName() valid = %v, want %v", valid, tt.wantValid)
			}

			if message != tt.wantMessage {
				t.Errorf("ValidateName() message = %v, want %v", message, tt.wantMessage)
			}
		})
	}
}

func TestValidateRegisterUserRequest(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		firstName  string
		lastName   string
		wantErrors map[string]string
	}{
		{
			name:       "All valid",
			email:      "test@example.com",
			password:   "Password123",
			firstName:  "John",
			lastName:   "Doe",
			wantErrors: map[string]string{},
		},
		{
			name:      "All invalid",
			email:     "invalid-email",
			password:  "short",
			firstName: "J",
			lastName:  "D",
			wantErrors: map[string]string{
				"email":      "Invalid email format",
				"password":   "Password must be at least 8 characters",
				"first_name": "First name must be at least 2 characters",
				"last_name":  "Last name must be at least 2 characters",
			},
		},
		{
			name:      "Invalid email only",
			email:     "invalid",
			password:  "Password123",
			firstName: "John",
			lastName:  "Doe",
			wantErrors: map[string]string{
				"email": "Invalid email format",
			},
		},
		{
			name:      "Invalid password only",
			email:     "test@example.com",
			password:  "short",
			firstName: "John",
			lastName:  "Doe",
			wantErrors: map[string]string{
				"password": "Password must be at least 8 characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateRegisterUserRequest(tt.email, tt.password, tt.firstName, tt.lastName)

			if len(errors) != len(tt.wantErrors) {
				t.Errorf("ValidateRegisterUserRequest() errors count = %v, want %v", len(errors), len(tt.wantErrors))
			}

			for key, wantMsg := range tt.wantErrors {
				if gotMsg, ok := errors[key]; !ok {
					t.Errorf("ValidateRegisterUserRequest() missing error for %v", key)
				} else if gotMsg != wantMsg {
					t.Errorf("ValidateRegisterUserRequest() error[%v] = %v, want %v", key, gotMsg, wantMsg)
				}
			}
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		wantErrors map[string]string
	}{
		{
			name:       "Valid login",
			email:      "test@example.com",
			password:   "anypassword",
			wantErrors: map[string]string{},
		},
		{
			name:     "Invalid email",
			email:    "invalid-email",
			password: "anypassword",
			wantErrors: map[string]string{
				"email": "Invalid email format",
			},
		},
		{
			name:     "Empty password",
			email:    "test@example.com",
			password: "",
			wantErrors: map[string]string{
				"password": "Password is required",
			},
		},
		{
			name:     "Both invalid",
			email:    "",
			password: "",
			wantErrors: map[string]string{
				"email":    "Email is required",
				"password": "Password is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateLoginRequest(tt.email, tt.password)

			if len(errors) != len(tt.wantErrors) {
				t.Errorf("ValidateLoginRequest() errors count = %v, want %v", len(errors), len(tt.wantErrors))
			}

			for key, wantMsg := range tt.wantErrors {
				if gotMsg, ok := errors[key]; !ok {
					t.Errorf("ValidateLoginRequest() missing error for %v", key)
				} else if gotMsg != wantMsg {
					t.Errorf("ValidateLoginRequest() error[%v] = %v, want %v", key, gotMsg, wantMsg)
				}
			}
		})
	}
}

func TestValidateUpdateUserRequest(t *testing.T) {
	validEmail := "test@example.com"
	invalidEmail := "invalid"
	validName := "John"
	shortName := "J"

	tests := []struct {
		name       string
		email      *string
		firstName  *string
		lastName   *string
		wantErrors map[string]string
	}{
		{
			name:       "All valid",
			email:      &validEmail,
			firstName:  &validName,
			lastName:   &validName,
			wantErrors: map[string]string{},
		},
		{
			name:       "Nil fields (no validation)",
			email:      nil,
			firstName:  nil,
			lastName:   nil,
			wantErrors: map[string]string{},
		},
		{
			name:      "Invalid email",
			email:     &invalidEmail,
			firstName: &validName,
			lastName:  &validName,
			wantErrors: map[string]string{
				"email": "Invalid email format",
			},
		},
		{
			name:      "Invalid first name",
			email:     &validEmail,
			firstName: &shortName,
			lastName:  &validName,
			wantErrors: map[string]string{
				"first_name": "First name must be at least 2 characters",
			},
		},
		{
			name:      "Invalid last name",
			email:     &validEmail,
			firstName: &validName,
			lastName:  &shortName,
			wantErrors: map[string]string{
				"last_name": "Last name must be at least 2 characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateUpdateUserRequest(tt.email, tt.firstName, tt.lastName)

			if len(errors) != len(tt.wantErrors) {
				t.Errorf("ValidateUpdateUserRequest() errors count = %v, want %v", len(errors), len(tt.wantErrors))
			}

			for key, wantMsg := range tt.wantErrors {
				if gotMsg, ok := errors[key]; !ok {
					t.Errorf("ValidateUpdateUserRequest() missing error for %v", key)
				} else if gotMsg != wantMsg {
					t.Errorf("ValidateUpdateUserRequest() error[%v] = %v, want %v", key, gotMsg, wantMsg)
				}
			}
		})
	}
}