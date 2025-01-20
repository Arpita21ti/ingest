package requests

// StudentSignUpRequest represents the unified model for user signup data.
type StudentSignUpRequest struct {
	// Common Fields for master lookup
	EnrollmentNo string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required,len=12,enrollmentNo"`

	// Student Login Credentials
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
	Phone    string `json:"phone" bson:"phone" validate:"required,phone"`

	// Student Academic Details
	Branch                string  `json:"branch" bson:"branch" validate:"required,branch"`
	YearOfEnrollment      int     `json:"yearOfEnrollment" bson:"yearOfEnrollment" validate:"required,yearOfEnrollment"`
	CGPA                  float32 `json:"cgpa" bson:"cgpa" validate:"gte=0,lte=10"`
	PreviousSemSGPA       float32 `json:"previousSemSgpa" bson:"previousSemSgpa" validate:"gte=0,lte=10"`
	SchoolForClassTen     string  `json:"schoolForClassTen" bson:"schoolForClassTen" validate:"required"`
	ClassTenPercentage    float32 `json:"classTenPercentage" bson:"classTenPercentage" validate:"gte=0,lte=100"`
	SchoolForClassTwelve  string  `json:"schoolForClassTwelve" bson:"schoolForClassTwelve" validate:"required"`
	ClassTwelvePercentage float32 `json:"classTwelvePercentage" bson:"classTwelvePercentage" validate:"gte=0,lte=100"`

	// TODO: Update
	// student Certifications and Achievements Details
	Certifications []string `json:"certifications" bson:"certifications" validate:"required"`

	IssuingAuthority   []string `json:"issuingAuthority" bson:"issuingAuthority" validate:"required"`
	CertificationNames []string `json:"certificationNames" bson:"certificationNames" validate:"required"`
	IssuingDate        []string `json:"issuingDate" bson:"issuingDate" validate:"required"`

	// Family Details
	FatherName          string `json:"fatherName" bson:"fatherName" validate:"required"`
	FatherQualification string `json:"fatherQualification" bson:"fatherQualification" validate:"required"`
	FatherProfession    string `json:"fatherProfession" bson:"fatherProfession" validate:"required"`
	MotherName          string `json:"motherName" bson:"motherName" validate:"required"`
	MotherQualification string `json:"motherQualification" bson:"motherQualification" validate:"required"`
	MotherProfession    string `json:"motherProfession" bson:"motherProfession" validate:"required"`
	NoOfSiblings        int    `json:"noOfSiblings" bson:"noOfSiblings" validate:"required"`
	TotalFamilyIncome   int    `json:"totalFamilyIncome" bson:"totalFamilyIncome" validate:"required"`

	// Student Profile Information
	Name     string `json:"name" bson:"name" validate:"required,name"`
	Gender   string `json:"gender" bson:"gender" validate:"required,oneof=M F O"`
	Category string `json:"category" bson:"category" validate:"required,oneof=GEN EWS OBC SC ST"`

	// Scholarship Details
	ScholarshipName string `json:"scholarshipName" bson:"scholarshipName" validate:"required"`
	ProvidedBy      string `json:"providedBy" bson:"providedBy" validate:"required"`
	AmountReceived  int    `json:"amountReceived" bson:"amountReceived" validate:"required"`
}
