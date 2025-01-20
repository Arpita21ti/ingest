package models

type StudentResetPasswordRequest struct {
	EnrollmentNo string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required"`
	NewPassword  string `json:"newPassword" binding:"required,min=8"`
}
