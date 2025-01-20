// DTO (Data Transfer Object) for the response of the GetNicheResponse API
package response

type GetNicheResponse struct {
	QuestionNicheID uint32 `json:"nicheID" bson:"nicheID"`
	NicheName       string `json:"nicheName" bson:"nicheName"`
	// Added bor backtrack and frontend filtering capabilities
	QuestionSubDomainID uint32 `json:"subDomainID" bson:"subDomainID"`
}
