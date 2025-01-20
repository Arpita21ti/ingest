package models

type StudentCollegeAcademics struct {
	EnrollmentNo    string  `json:"enrollmentNo" bson:"enrollmentNo" binding:"required"`
	PreviousSemCGPA float32 `json:"previousSemCGPA" bson:"previousSemCGPA" binding:"required,gte=0,lte=10"` // Validating CGPA between 0 and 10
	PreviousSemSGPA float32 `json:"previousSemSGPA" bson:"previousSemSGPA" binding:"required,gte=0,lte=10"` // Validating SGPA between 0 and 10
}
