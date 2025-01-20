package models

type TextBasedQuestion struct {
	BaseQuestion // Embedding common fields
}

func (TextBasedQuestion) TableName() string {
	return "question_schema.text_questions"
}
