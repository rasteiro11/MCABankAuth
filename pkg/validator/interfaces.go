package validator

//go:generate mockery --name=EmailValidator --output=./mocks --case=underscore
type EmailValidator interface {
	IsValid(email string) bool
}
