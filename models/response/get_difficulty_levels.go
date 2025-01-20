// DTO (Data Transfer Object) for the response of the GetDifficultyLevelsResponse API
package response

type GetDifficultyLevelsResponse struct {
	QuestionDifficultyLevelID uint32 `json:"difficultyLevelID" bson:"difficultyLevelID"`
	DifficultyLevel           string `json:"difficultyLevel" bson:"difficultyLevel"`
	// Added bor backtrack and frontend filtering capabilities
	QuestionNicheID uint32 `json:"nicheID" bson:"nicheID"`
}
