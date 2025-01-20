package response

type GetQuestionsResponse struct {
	Questions         interface{} `json:"questions"`
	PracticeSessionID uint32      `json:"practiceSessionID"`
	Message           string      `json:"message"`
}
