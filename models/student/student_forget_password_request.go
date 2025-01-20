package models

type StudentForgotPasswordRequest struct {
	EnrollmentNo string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	Phone        string `json:"phone" bson:"phone" validate:"required,phone"`
}
