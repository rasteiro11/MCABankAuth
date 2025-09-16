package security

//go:generate mockery --name=PasswordHasher --output=./mocks --case=underscore
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}
