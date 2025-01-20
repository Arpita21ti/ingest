// DTO (Data Transfer Object) for the response of the GetDomains API
package response

type GetDomainsResponse struct {
	QuestionDomainID uint32 `json:"domainID" bson:"domainID"`
	DomainName       string `json:"domainName" bson:"domainName"`
}
