package requests

type SuccessfullyEndPracticeSessionRequest struct {
	// PracticeSessionID = Unique identifier for the practice session to end
	PracticeSessionID uint32 `json:"practiceSessionId" binding:"required"`
	// QuestionsAttempted = Number of questions attempted during the session
	QuestionsAttempted int `json:"questionsAttempted" binding:"required"`
	// QuestionsCorrect = Number of questions answered correctly
	QuestionsCorrect int `json:"questionsCorrect" binding:"required"`
	// Feedbacks = Feedbacks given by the student for the practice session
	Feedbacks string `json:"feedbacks" binding:"required"`

	// ScoreEarnedPercentage = Total score percentage ((correct ans/total ques) * 100) earned in the session
	// Calculate and send from the front end.
	ScoreEarnedPercentage float64 `json:"scoreEarnedPercentage" binding:"required"`
}
