package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailValidator(t *testing.T) {
	v := NewEmailValidator()

	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"john.doe@mail.co.uk", true},
		{"invalid-email", false},
		{"missing@domain", false},
		{"@missinguser.com", false},
		{"test@.com", false},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, v.IsValid(tt.email), "email: %s", tt.email)
	}
}
