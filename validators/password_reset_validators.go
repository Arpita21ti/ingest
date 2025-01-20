package validators

import "github.com/go-playground/validator/v10"

// RegisterValidators registers custom validation functions
func RegisterValidatorsPasswordReset(validate *validator.Validate) {
	validate.RegisterValidation("phone", phoneValidator)
}
