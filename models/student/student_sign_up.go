package models

// StudentSignUp represents the unified model for user signup data.
type StudentSignUp struct {
	// User Login Credentials
	EnrollmentNo string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	Password     string `json:"password" bson:"password" validate:"required"`
	Phone        string `json:"phone" bson:"phone" validate:"required,phone"`

	// Student Profile Information
	Name               string  `json:"name" bson:"name" validate:"required,name"`
	Branch             string  `json:"branch" bson:"branch" validate:"required,branch"`
	YearOfEnrollment   int     `json:"yearOfEnrollment" bson:"yearOfEnrollment" validate:"required,yearOfEnrollment"`
	ClassTenPercent    float32 `json:"classTenPercent" bson:"classTenPercent" validate:"required,percentage"`
	ClassTwelvePercent float32 `json:"classTwelvePercent" bson:"classTwelvePercent" validate:"required,percentage"`

	// Student Academic Information
	PreviousSemCGPA float32 `json:"previousSemCGPA" bson:"previousSemCGPA" validate:"required,gte=0,lte=10,cgpa"` // Validating CGPA between 0 and 10
	PreviousSemSGPA float32 `json:"previousSemSGPA" bson:"previousSemSGPA" validate:"required,gte=0,lte=10,cgpa"`

	// Student Certifications and Resume
	Certifications []string `json:"certifications" bson:"certifications" validate:"required,certificate"`
	Resume         string   `json:"resume" bson:"resume" binding:"required,url" validate:"url"`
}
