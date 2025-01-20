package models

type QuestionDifficultyLevelTable struct {
	QuestionDifficultyLevelID uint32 `gorm:"primaryKey;autoIncrement;unique" json:"difficultyID" bson:"difficultyID"`
	QuestionNicheID           uint32 `gorm:"not null;" json:"nicheID" bson:"nicheID"`
	DifficultyLevel           string `gorm:"type:varchar(6);not null;check:difficulty_level IN ('EASY','MEDIUM','HARD')" json:"difficultyLevel" bson:"difficultyLevel" binding:"required,oneof=EASY MEDIUM HARD"`

	// Foreign key for QuestionFormatTable table.
	// The foreign key will be made in the QuestionFormatTable table.
	// Relationships
	Difficulty QuestionFormatTable `gorm:"foreignKey:question_difficulty_level_id;references:question_difficulty_level_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (QuestionDifficultyLevelTable) TableName() string {
	return "question_schema.question_difficulty_level_table"
}
