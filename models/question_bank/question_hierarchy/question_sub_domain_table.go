package models

type QuestionSubDomainsTable struct {
	QuestionSubDomainID uint32 `gorm:"primaryKey;autoIncrement;unique" json:"subDomainID" bson:"subDomainID"`
	QuestionDomainID    uint32 `gorm:"not null;" json:"domainID" bson:"domainID"`
	SubDomainName       string `gorm:"type:varchar(255);not null;" json:"subDomainName" bson:"subDomainName"`

	// Foreign key for QuestionNicheTable table.
	// The foreign key will be made in the QuestionNicheTable table.
	// Relationships
	SubDomain QuestionNicheTable `gorm:"foreignKey:question_sub_domain_id;references:question_sub_domain_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (QuestionSubDomainsTable) TableName() string {
	return "question_schema.question_sub_domains_table"
}
