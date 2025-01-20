package models

type StudentProfile struct {
	EnrollmentNo       string  `json:"enrollmentNo" bson:"enrollmentNo" binding:"required"`
	Name               string  `json:"name" bson:"name" binding:"required"`
	Branch             string  `json:"branch" bson:"branch" binding:"required"`
	YearOfEnrollment   int     `json:"yearOfEnrollment" bson:"yearOfEnrollment" binding:"required"`
	ClassTenPercent    float32 `json:"classTenPercent" bson:"classTenPercent" binding:"required"`
	ClassTwelvePercent float32 `json:"classTwelvePercent" bson:"classTwelvePercent" binding:"required"`
}
