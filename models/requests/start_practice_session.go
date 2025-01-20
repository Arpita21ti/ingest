package requests

type StartPracticeSessionRequest struct {
	// EnrollmentNo = Unique identifier for each student
	EnrollmentNo string `json:"enrollmentNo" binding:"required"`

	// Domain = Category of the questions (e.g., Programming, Mathematics)
	DomainID uint32 `json:"domainID" binding:"required"`

	// SubDomain = Sub-category of the questions (e.g., Data Structures, Algebra)
	SubDomainID uint32 `json:"subDomainID" binding:"required"`

	// Difficulty = Difficulty level of the session (e.g., Easy, Medium, Hard)
	DifficultyLevelID uint32 `json:"difficultyLevelID" binding:"required"`
}
