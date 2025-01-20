package requests

type StudentLoginRequest struct {
	EnrollmentNo string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	Password     string `json:"password" bson:"password" validate:"required"`
	Phone        string `json:"phone" bson:"phone" validate:"required,phone"`
}
