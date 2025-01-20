// DTO (Data Transfer Object) for the response of the GetQuestionFormatsResponse API
package response

type GetQuestionFormatsResponse struct {
	QuestionFormatID uint32 `json:"formatID" bson:"formatID"`
	Format           string `json:"format" bson:"format"`
	// Added bor backtrack and frontend filtering capabilities
	QuestionDifficultyLevelID uint32 `json:"difficultyID" bson:"difficultyID"`
}
