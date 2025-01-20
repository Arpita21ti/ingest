package validators

import "github.com/go-playground/validator/v10"

// Register custom validation functions for Practice Session
func RegisterValidatorsPracticeSession(validate *validator.Validate) {
	validate.RegisterValidation("enrollmentNo", enrollmentNoValidator)
}
