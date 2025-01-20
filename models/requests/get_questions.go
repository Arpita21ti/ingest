package requests

type GetQuestionsRequest struct {
	EnrollmentNo              string `json:"enrollmentNo" bson:"enrollmentNo" validate:"required,enrollmentNo"`
	QuestionDomainID          uint32 `json:"questionDomainID" bson:"questionDomainID" binding:"required"`
	QuestionSubDomainID       uint32 `json:"questionSubDomainID" bson:"questionSubDomainID" binding:"required"`
	QuestionDifficultyLevelID uint32 `json:"questionDifficultyLevelID" bson:"questionDifficultyLevelID" binding:"required"`
	QuestionFormatID          uint32 `json:"questionFormatID" bson:"questionFormatID" binding:"required"`
	QuestionFormat            string `json:"questionFormat" bson:"questionFormat" binding:"required,max=3"`
	QuestionCount             int    `json:"questionCount" bson:"questionCount" binding:"required"`
	LastAttemptedQuestionID   uint32 `json:"lastAttemptedQuestionID" bson:"lastAttemptedQuestionID" binding:"required"`
}
