package models

import "time"

// PasswordResetRequest represents a password reset request document.
type StudentPasswordResetRequestsStorage struct {
	EnrollmentNo     string    `bson:"enrollmentNo" json:"enrollmentNo"`
	ResetToken       string    `bson:"reset_token" json:"reset_token"`
	ResetTokenExpiry time.Time `bson:"reset_token_expiry" json:"reset_token_expiry"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}
