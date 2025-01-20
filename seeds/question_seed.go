package seed

import (
	"log"
	"server/config"
	questionbank "server/models/question_bank/question_type"

	"github.com/lib/pq"
)

// SeedQuestions seeds the database with test questions
func SeedQuestions() error {
	db := config.GetPostgresDBConnection()

	// Create sample MCQ questions
	mcqQuestions := []questionbank.MCQQuestion{
		{
			BaseQuestion: questionbank.BaseQuestion{
				// Domain:        "Aptitude",
				// Subcategories: pq.StringArray{"Probability", "Linear Algebra"}, // Correct array format
				QuestionText: "What is the probability of getting a 6 on a fair die?",
				Answer:       "1/6",
				// Difficulty:    "Medium",
				// Tags:          pq.StringArray{"basic", "probability"},
			},
			Explanation: "A standard die has 6 faces, so the probability of getting a 6 is 1 out of 6.",
			Options:     pq.StringArray{"1/6", "1/4", "1/3", "1/5"},
		},
	}

	// Insert MCQ questions into the database
	for _, question := range mcqQuestions {
		if err := db.Create(&question).Error; err != nil {
			log.Fatalf("Error inserting MCQ question: %v", err)
			return err // Return the error if insertion fails
		}
	}

	// Create sample True/False questions
	tfQuestions := []questionbank.TrueFalseQuestion{
		{
			BaseQuestion: questionbank.BaseQuestion{
				// Domain:        "Aptitude",
				// Subcategories: pq.StringArray{"Critical Thinking"},
				QuestionText: "The Earth is flat.",
				Answer:       "False",
				// Difficulty:    "Easy",
				// Tags:          pq.StringArray{"science", "geography"},
			},
		},
	}

	// Insert True/False questions into the database
	for _, question := range tfQuestions {
		if err := db.Create(&question).Error; err != nil {
			log.Fatalf("Error inserting True/False question: %v", err)
			return err // Return the error if insertion fails
		}
	}

	// Create sample Fill-in-the-Blank questions
	fbQuestions := []questionbank.FillInTheBlankQuestion{
		{
			BaseQuestion: questionbank.BaseQuestion{
				// Domain:        "General Knowledge",
				// Subcategories: pq.StringArray{"Geography"},
				QuestionText: "The capital of France is ________.",
				Answer:       "Paris",
				// Difficulty:    "Medium",
				// Tags:          pq.StringArray{"geography", "cities"},
			},
		},
	}

	// Insert Fill-in-the-Blank questions into the database
	for _, question := range fbQuestions {
		if err := db.Create(&question).Error; err != nil {
			log.Fatalf("Error inserting Fill-in-the-Blank question: %v", err)
			return err // Return the error if insertion fails
		}
	}

	// Create sample Text-based questions
	txtQuestions := []questionbank.TextBasedQuestion{
		{
			BaseQuestion: questionbank.BaseQuestion{
				// Domain:        "General Knowledge",
				// Subcategories: pq.StringArray{"History"},
				QuestionText: "Explain the significance of the French Revolution.",
				Answer:       "It led to the rise of democracy and the fall of absolute monarchies.",
				// Difficulty:    "Hard",
				// Tags:          pq.StringArray{"history", "revolution"},
			},
		},
	}

	// Insert Text-based questions into the database
	for _, question := range txtQuestions {
		if err := db.Create(&question).Error; err != nil {
			log.Fatalf("Error inserting Text-based question: %v", err)
			return err // Return the error if insertion fails
		}
	}

	// If everything was successful, log a success message and return nil
	log.Println("Test questions inserted successfully.")
	return nil
}
