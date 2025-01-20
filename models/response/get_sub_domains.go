// DTO (Data Transfer Object) for the response of the GetSubDomains API
package response

type GetSubDomainsResponse struct {
	QuestionSubDomainID uint32 `json:"subDomainID" bson:"subDomainID"`
	SubDomainName       string `json:"subDomainName" bson:"subDomainName"`
	// Added bor backtrack and frontend filtering capabilities
	QuestionDomainID uint32 `json:"domainID" bson:"domainID"`
}
