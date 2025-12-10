package validator

import (
	"fmt"

	"github.com/Infamous003/follow-service/internal/domain"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Errors)
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(ok bool, key string, value string) {
	if !ok {
		v.AddError(key, value)
	}
}

func ValidateUser(v *Validator, user *domain.User) {
	v.Check(user.Username != "", "username", "Username is required")
	v.Check(len(user.Username) <= 32, "username", "Username must not be longer than 32 characters")

	// there can be more validations added here later
}
