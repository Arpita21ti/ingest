package models

type TrueFalseQuestion struct {
	BaseQuestion        // Embedding common fields
	Explanation  string `gorm:"type:text;default:''" json:"explanation" bson:"explanation"`
}

func (TrueFalseQuestion) TableName() string {
	return "question_schema.tf_questions"
}
