package validators

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// enrollmentNoValidator validates the enrollment number format
func enrollmentNoValidator(fl validator.FieldLevel) bool {
	enrollmentNo := fl.Field().String()

	// Regular expression to match the format: 4 digits, 2 letters, 6 digits (e.g., 1234AB567890)
	re := regexp.MustCompile(`^\d{4}[A-Za-z]{2}\d{6}$`)
	return re.MatchString(enrollmentNo)
}

// phoneValidator validates phone numbers
func phoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return re.MatchString(phone)
}

// nameValidator ensures the name contains only alphabets and spaces, and has a valid length
func nameValidator(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	// Regex to allow only alphabets and spaces, with 2 to 50 characters
	re := regexp.MustCompile(`^[A-Za-z\s]{2,50}$`)
	return re.MatchString(name)
}

// branchValidator ensures the branch matches a predefined list
func branchValidator(fl validator.FieldLevel) bool {
	branch := strings.ToUpper(fl.Field().String()) // Normalize to uppercase

	// List of allowed Branch Codes.
	validBranches := map[string]bool{
		"AIML":  true,
		"CSDS":  true,
		"CSBS":  true,
		"CSE":   true,
		"ECE":   true,
		"ME":    true,
		"EEE":   true,
		"IT":    true,
		"CIVIL": true,
		"CHEM":  true,
		// "BIO": true,
		// "AERONAUTICAL": true,
		"BIOTECH": true,
		// "MATERIALS": true,
		// "ENVIRONMENTAL": true,
		"PCT": true,
	}

	_, exists := validBranches[branch]
	return exists
}

// yearOfEnrollmentValidator checks if the year is within a valid range
func yearOfEnrollmentValidator(fl validator.FieldLevel) bool {
	year, ok := fl.Field().Interface().(int) // Ensure the field is an integer
	fmt.Println(ok, year)
	if !ok {
		return false // Not an integer
	}

	currentYear := time.Now().Year()           // Get the current year
	return year >= 1990 && year <= currentYear // Allowed range of 1900 to present date.
}

// percentageValidator validates percentage between 0 and 100
func percentageValidator(fl validator.FieldLevel) bool {
	percent := fl.Field().Float()
	return percent >= 0 && percent <= 100
}

// gPAValidator validates CGPA between 0 and 10
func gPAValidator(fl validator.FieldLevel) bool {
	gpa := fl.Field().Float()
	return gpa >= 0 && gpa <= 10
}

// certificationsValidator validates that every element in a slice is a valid URL
func certificationsValidator(fl validator.FieldLevel) bool {
	certifications, ok := fl.Field().Interface().([]string)
	if !ok {
		return false // Not a slice of strings
	}

	// Not in loop to avoid re-calculation every time.
	re := regexp.MustCompile(`^(http|https)://[^\s/$.?#].[^\s]*$`)

	for _, cert := range certifications {
		if !re.MatchString(cert) {
			return false // Invalid URL found
		}
	}
	return true // All URLs are valid
}

// uRLValidator validates if the given string is a valid URL
func uRLValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(http|https)://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(fl.Field().String())
}

// TODO: Add more validators for signup fields.
// RegisterValidators registers custom validation functions
func RegisterValidatorsSignUp(validate *validator.Validate) {
	validate.RegisterValidation("enrollmentNo", enrollmentNoValidator)
	validate.RegisterValidation("phone", phoneValidator)
	validate.RegisterValidation("percentage", percentageValidator)
	validate.RegisterValidation("gpa", gPAValidator)
	validate.RegisterValidation("url", uRLValidator)
	validate.RegisterValidation("certificate", certificationsValidator)
	validate.RegisterValidation("yearOfEnrollment", yearOfEnrollmentValidator)
	validate.RegisterValidation("name", nameValidator)
	validate.RegisterValidation("branch", branchValidator)
}

// RegisterValidators registers custom validation functions
func RegisterValidatorsLogIn(validate *validator.Validate) {
	validate.RegisterValidation("phone", phoneValidator)
}
