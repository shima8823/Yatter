package object_test

import (
	"testing"
	"yatter-backend-go/app/domain/object"
)

func TestSetPassword(t *testing.T) {
	acc := &object.Account{}
	password := "secret_password"

	err := acc.SetPassword(password)
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	if acc.PasswordHash == "" {
		t.Fatalf("expected password hash to be set, but got an empty hash")
	}
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name          string
		accountPass   string
		inputPass     string
		expectedMatch bool
	}{
		{
			name:          "Matching Password",
			accountPass:   "test123",
			inputPass:     "test123",
			expectedMatch: true,
		},
		{
			name:          "Non-Matching Password",
			accountPass:   "test123",
			inputPass:     "test456",
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &object.Account{}
			err := acc.SetPassword(tt.accountPass)
			if err != nil {
				t.Fatalf("setting password failed: %v", err)
			}

			match := acc.CheckPassword(tt.inputPass)
			if match != tt.expectedMatch {
				t.Fatalf("expected match = %v, but got %v", tt.expectedMatch, match)
			}
		})
	}
}
