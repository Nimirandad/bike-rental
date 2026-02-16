package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) (bool, string) {
	email = strings.TrimSpace(email)
	if email == "" {
		return false, "Email is required"
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false, "Invalid email format"
	}

	return true, ""
}

func ValidatePassword(password string) (bool, string) {
	if password == "" {
		return false, "Password is required"
	}

	if len(password) < 8 {
		return false, "Password must be at least 8 characters"
	}

	if len(password) > 100 {
		return false, "Password is too long (max 100 characters)"
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return false, "Password must contain at least one letter and one number"
	}

	return true, ""
}

func ValidateName(name, fieldName string) (bool, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, fieldName + " is required"
	}

	if len(name) < 2 {
		return false, fieldName + " must be at least 2 characters"
	}

	if len(name) > 50 {
		return false, fieldName + " is too long (max 50 characters)"
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-]+$`)
	if !nameRegex.MatchString(name) {
		return false, fieldName + " can only contain letters, spaces, and hyphens"
	}

	return true, ""
}

func ValidateRegisterUserRequest(email, password, firstName, lastName string) map[string]string {
	errors := make(map[string]string)

	if valid, msg := ValidateEmail(email); !valid {
		errors["email"] = msg
	}

	if valid, msg := ValidatePassword(password); !valid {
		errors["password"] = msg
	}

	if valid, msg := ValidateName(firstName, "First name"); !valid {
		errors["first_name"] = msg
	}

	if valid, msg := ValidateName(lastName, "Last name"); !valid {
		errors["last_name"] = msg
	}

	return errors
}

func ValidateLoginRequest(email, password string) map[string]string {
	errors := make(map[string]string)

	if valid, msg := ValidateEmail(email); !valid {
		errors["email"] = msg
	}

	if password == "" {
		errors["password"] = "Password is required"
	}

	return errors
}

func ValidateUpdateUserRequest(email, firstName, lastName *string) map[string]string {
	errors := make(map[string]string)

	if email != nil && *email != "" {
		if valid, msg := ValidateEmail(*email); !valid {
			errors["email"] = msg
		}
	}

	if firstName != nil && *firstName != "" {
		if valid, msg := ValidateName(*firstName, "First name"); !valid {
			errors["first_name"] = msg
		}
	}

	if lastName != nil && *lastName != "" {
		if valid, msg := ValidateName(*lastName, "Last name"); !valid {
			errors["last_name"] = msg
		}
	}

	return errors
}