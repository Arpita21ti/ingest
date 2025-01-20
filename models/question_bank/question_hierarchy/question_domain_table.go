package models

type QuestionDomainsTable struct {
	QuestionDomainID uint32 `gorm:"primaryKey;autoIncrement;unique" json:"domainId" bson:"domainId"`
	DomainName       string `gorm:"type:varchar(255);not null" json:"domainName" bson:"domainName"`

	// Foreign key for QuestionSubDomainsTable table.
	// The foreign key will be made in the QuestionSubDomainsTable table.
	// Relationships
	Domain QuestionSubDomainsTable `gorm:"foreignKey:question_domain_id;refernces:question_domain_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-" bson:"-"`
}

func (QuestionDomainsTable) TableName() string {
	return "question_schema.question_domains_table"
}
