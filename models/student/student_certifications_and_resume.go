package models

type StudentCertificationsAndResume struct {
	EnrollmentNo   string   `json:"enrollmentNo" bson:"enrollmentNo" binding:"required"`
	Certifications []string `json:"certifications" bson:"certifications" binding:"required"`
	Resume         string   `json:"resume" bson:"resume" binding:"required,url"`
}
