package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHasher(t *testing.T) {
	hasher := NewPasswordHasher()

	t.Run("Hash and Verify - success", func(t *testing.T) {
		password := "SuperSecret123!"
		hashed, err := hasher.Hash(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)

		assert.True(t, hasher.Verify(password, hashed))
	})

	t.Run("Verify - failure", func(t *testing.T) {
		password := "SuperSecret123!"
		hashed, err := hasher.Hash(password)
		assert.NoError(t, err)

		assert.False(t, hasher.Verify("WrongPassword!", hashed))
	})
}
