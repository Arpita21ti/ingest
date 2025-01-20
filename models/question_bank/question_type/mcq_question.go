package models

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// MCQQuestion extends BaseQuestion for MCQ questions
type MCQQuestion struct {
	BaseQuestion                // Embedding common fields
	Explanation  string         `gorm:"type:text;default:''"  json:"explanation,omitempty" bson:"explanation,omitempty"`
	Options      pq.StringArray `gorm:"not null;type:text[];not null" json:"options" bson:"options"` // Additional field for MCQ options
}

func (MCQQuestion) TableName() string {
	return "question_schema.mcq_questions"
}

// Use if want to check before every save operation of the database.
// BeforeSave is a GORM hook that validates the MCQ question before saving
func (m *MCQQuestion) BeforeSave(tx *gorm.DB) error {
	if len(m.Options) != 4 {
		return errors.New("MCQ questions must have exactly 4 options")
	}
	return nil
}

// Use if want to use custom or multilevel validations and want to call check and validations
// conditionally and not every time before inserting or updating data.
func (m *MCQQuestion) ValidateOptions() error {
	if len(m.Options) != 4 {
		return errors.New("MCQ questions must have exactly 4 options")
	}
	return nil
}
