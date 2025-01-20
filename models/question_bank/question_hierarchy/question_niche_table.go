package models

type QuestionNicheTable struct {
	QuestionNicheID     uint32 `gorm:"primaryKey;autoIncrement" json:"nicheID" bson:"nicheID"`
	QuestionSubDomainID uint32 `gorm:"not null;" json:"subDomainID" bson:"subDomainID"`
	NicheName           string `gorm:"type:varchar(255);not null" json:"nicheName" bson:"nicheName"`

	// Foreign key for QuestionDifficultyLevelTable table.
	// The foreign key will be made in the QuestionDifficultyLevelTable table.
	// Relationships
	Niche QuestionDifficultyLevelTable `gorm:"foreignKey:QuestionNicheID;references:QuestionNicheID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (QuestionNicheTable) TableName() string {
	return "question_schema.question_niches_table"
}
