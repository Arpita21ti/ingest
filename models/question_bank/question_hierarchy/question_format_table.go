package models

import (
	question_type "server/models/question_bank/question_type"
)

type QuestionFormatTable struct {
	QuestionFormatID          uint32 `gorm:"primaryKey;autoIncrement;unique" json:"formatID" bson:"formatID"`
	QuestionDifficultyLevelID uint32 `gorm:"not null;" json:"difficultyID" bson:"difficultyID"`
	Format                    string `gorm:"type:varchar(3);not null;check:format IN('TXT','MCQ','FIB','TF')" json:"format" bson:"format" binding:"required,oneof=TXT MCQ FIB TF"`

	// Relationships
	QuestionFormat question_type.BaseQuestion `gorm:"foreignKey:question_format_id;references:question_format_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (QuestionFormatTable) TableName() string {
	return "question_schema.question_formats_table"
}
