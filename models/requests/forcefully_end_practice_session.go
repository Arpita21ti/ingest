package requests

type ForcefullyEndPracticeSessionRequest struct {
	// PracticeSessionID = Unique identifier for the practice session to end
	PracticeSessionID uint32 `json:"practiceSessionId" binding:"required"`
	// EnrollmentNo = Unique identifier for each student
	EnrollmentNo string `json:"enrollmentNo" binding:"required"`
}
