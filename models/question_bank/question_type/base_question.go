package models

import (
	"time"
)

// BaseQuestion holds fields common to all question models.
type BaseQuestion struct {
	QuestionFormatID uint32    `gorm:"not null;primaryKey" json:"formatID" bson:"formatID"`          // Part of composite primary key
	QuestionID       uint32    `gorm:"autoIncrement;primaryKey" json:"questionID" bson:"questionID"` // Part of composite primary key
	QuestionText     string    `gorm:"type:text;not null" json:"questionText" bson:"questionText"`
	Answer           string    `gorm:"type:text;not null" json:"answer" bson:"answer"`
	UpdatedAt        time.Time `gorm:"not null;autoUpdateTime;index"`
}
