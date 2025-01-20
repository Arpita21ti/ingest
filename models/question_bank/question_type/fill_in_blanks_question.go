package models

// FillInTheBlankQuestion extends BaseQuestion for fill-in-the-blank questions
type FillInTheBlankQuestion struct {
	BaseQuestion         // Embedding common fields
	Explanation  *string `gorm:"type:text;default:''" json:"explanation,omitempty" bson:"explanation,omitempty"`
}

func (FillInTheBlankQuestion) TableName() string {
	return "question_schema.fib_questions"
}
