package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"server/config"

	// questionbank "server/models/question_bank"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// FetchQuestionsHandler handles the hybrid fetch request
func FetchQuestionsHandler(c *gin.Context) {
	// Get category, difficulty, and type from the URL parameters
	category := c.Param("category")
	difficulty := c.DefaultQuery("difficulty", "") // Default to empty string if not provided
	questionType := c.DefaultQuery("type", "TXT")  // Default to TXT if not provided

	// Get the filters from the request body (other than category, difficulty, and type)
	var filters struct {
		Subcategories []string `json:"subcategories"`
		Tags          []string `json:"tags"`
	}
	if err := c.BindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters format"})
		return
	}

	// Log the received filters (for debugging purposes)
	fmt.Println("Category:", category)
	fmt.Println("Difficulty:", difficulty)
	fmt.Println("Subcategories:", filters.Subcategories)
	fmt.Println("Tags:", filters.Tags)
	fmt.Println("Type:", questionType)

	db := config.GetPostgresDBConnection()

	var questions interface{}

	// Create a base query for filtering based on category
	query := db.Where("category = ?", category)

	// Apply filters based on the URL parameters and body
	if len(filters.Tags) > 0 {
		query = query.Where("tags && ?", pq.Array(filters.Tags)) // PostgreSQL array overlap operator
	}
	if len(filters.Subcategories) > 0 {
		query = query.Where("subcategories && ?", pq.Array(filters.Subcategories)) // Array overlap for subcategories
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	// switch questionType {
	// case "TXT":
	// 	var txtQuestions []questionbank.TextBasedQuestion
	// 	if err := query.Model(&questionbank.TextBasedQuestion{}).Find(&txtQuestions).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch TXT questions", "details": err.Error()})
	// 		return
	// 	}
	// 	questions = txtQuestions
	// case "MCQ":
	// 	var mcqQuestions []questionbank.MCQQuestion
	// 	if err := query.Model(&questionbank.MCQQuestion{}).Find(&mcqQuestions).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch MCQ questions", "details": err.Error()})
	// 		return
	// 	}
	// 	questions = mcqQuestions
	// case "TF":
	// 	var tfQuestions []questionbank.TrueFalseQuestion
	// 	if err := query.Model(&questionbank.TrueFalseQuestion{}).Find(&tfQuestions).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch True/False questions", "details": err.Error()})
	// 		return
	// 	}
	// 	questions = tfQuestions
	// case "FB":
	// 	var fbQuestions []questionbank.FillInTheBlankQuestion
	// 	if err := query.Model(&questionbank.FillInTheBlankQuestion{}).Find(&fbQuestions).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Fill-in-the-Blank questions", "details": err.Error()})
	// 		return
	// 	}
	// 	questions = fbQuestions
	// default:
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question type"})
	// 	return
	// }

	// // Check if no questions were found
	// switch q := questions.(type) {
	// case []questionbank.TextBasedQuestion:
	// 	if len(q) == 0 {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found with the provided filters"})
	// 		return
	// 	}
	// case []questionbank.MCQQuestion:
	// 	if len(q) == 0 {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found with the provided filters"})
	// 		return
	// 	}
	// case []questionbank.TrueFalseQuestion:
	// 	if len(q) == 0 {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found with the provided filters"})
	// 		return
	// 	}
	// case []questionbank.FillInTheBlankQuestion:
	// 	if len(q) == 0 {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found with the provided filters"})
	// 		return
	// 	}
	// default:
	// 	log.Println("Unknown question type")
	// }

	// Return the filtered questions
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// AddSingleQuestionHandler handles the addition of a new question to the database
func AddSingleQuestionHandler(c *gin.Context) {
	// Get question type from the query parameter
	questionType := c.DefaultQuery("type", "TXT") // Default to TXT if not provided

	// db := config.GetPostgresDB()

	// Handle the question addition based on the question type
	switch questionType {
	// case "TXT":
	// 	var txtQuestion questionbank.TextBasedQuestion
	// 	if err := c.ShouldBindJSON(&txtQuestion); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TXT question format", "details": err.Error()})
	// 		return
	// 	}
	// 	if err := db.Create(&txtQuestion).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add TXT question", "details": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusCreated, gin.H{"message": "TXT question added successfully", "question": txtQuestion})

	// case "MCQ":
	// 	var mcqQuestion questionbank.MCQQuestion
	// 	if err := c.ShouldBindJSON(&mcqQuestion); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MCQ question format", "details": err.Error()})
	// 		return
	// 	}
	// 	if err := db.Create(&mcqQuestion).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add MCQ question", "details": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusCreated, gin.H{"message": "MCQ question added successfully", "question": mcqQuestion})

	// case "TF":
	// 	var tfQuestion questionbank.TrueFalseQuestion
	// 	if err := c.ShouldBindJSON(&tfQuestion); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid True/False question format", "details": err.Error()})
	// 		return
	// 	}
	// 	if err := db.Create(&tfQuestion).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add True/False question", "details": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusCreated, gin.H{"message": "True/False question added successfully", "question": tfQuestion})

	// case "FB":
	// 	var fbQuestion questionbank.FillInTheBlankQuestion
	// 	if err := c.ShouldBindJSON(&fbQuestion); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Fill-in-the-Blank question format", "details": err.Error()})
	// 		return
	// 	}
	// 	if err := db.Create(&fbQuestion).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Fill-in-the-Blank question", "details": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusCreated, gin.H{"message": "Fill-in-the-Blank question added successfully", "question": fbQuestion})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question type"})
		return
	}
}

func parseQuestions(input string) ([]interface{}, map[int]string, error) {
	var questions []interface{}              // To hold successfully parsed questions
	skippedQuestions := make(map[int]string) // Map to hold skipped questions with reasons
	entries := strings.Split(input, "---")   // Split input into individual question entries

	for i, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		lines := strings.Split(entry, "\n")
		questionType := ""
		// var question interface{}
		commonFields := make(map[string]bool) // Track presence of common fields
		var specializedFields bool            // Track specialized fields for the question type
		// var baseQuestion questionbank.BaseQuestion // The common base question for all specialized questions
		missingFields := []string{} // Collect missing fields
		skipEntry := false          // Flag to skip the current entry

		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Detect question type
			if strings.HasPrefix(line, "### Question Type:") {
				questionType = strings.TrimSpace(strings.TrimPrefix(line, "### Question Type:"))
				switch questionType {
				// case "MCQ":
				// 	question = &questionbank.MCQQuestion{}
				// case "FB":
				// 	question = &questionbank.FillInTheBlankQuestion{}
				// case "TXT":
				// 	question = &questionbank.TextBasedQuestion{}
				// case "TF":
				// 	question = &questionbank.TrueFalseQuestion{}
				default:
					missingFields = append(missingFields, "unknown question type")
					skipEntry = true // Mark this entry to skip
				}
			}

			if !skipEntry { // Process and store the common fields
				// if strings.HasPrefix(line, "Category:") {
				// 	baseQuestion.Domain = strings.TrimSpace(strings.TrimPrefix(line, "Category:"))
				// 	commonFields["Category"] = true
				// } else if strings.HasPrefix(line, "Subcategories:") {
				// 	baseQuestion.Subcategories = strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "Subcategories:")), ",")
				// 	commonFields["Subcategories"] = true
				// } else if strings.HasPrefix(line, "Question:") {
				// 	baseQuestion.QuestionText = strings.TrimSpace(strings.TrimPrefix(line, "Question:"))
				// 	commonFields["QuestionText"] = true
				// } else if strings.HasPrefix(line, "Answer:") {
				// 	baseQuestion.Answer = strings.TrimSpace(strings.TrimPrefix(line, "Answer:"))
				// 	commonFields["Answer"] = true
				// } else if strings.HasPrefix(line, "Difficulty:") {
				// 	baseQuestion.Difficulty = strings.TrimSpace(strings.TrimPrefix(line, "Difficulty:"))
				// 	commonFields["Difficulty"] = true
				// } else if strings.HasPrefix(line, "Tags:") {
				// 	baseQuestion.Tags = strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "Tags:")), ",")
				// 	commonFields["Tags"] = true
				// } else if strings.HasPrefix(line, "Status:") {
				// 	baseQuestion.Status = strings.TrimSpace(strings.TrimPrefix(line, "Status:"))
				// 	commonFields["Status"] = true
				// }

				// NOTE: Not stored the Question id, type, updatedAt  as it is set by default in the model.

				// Populate question-specific fields dynamically
				// switch q := question.(type) {
				// case *questionbank.MCQQuestion:
				// 	q.BaseQuestion = baseQuestion // Embed BaseQuestion

				// 	if strings.HasPrefix(line, "Option:") {
				// 		// Remove "Option:" prefix and trim any extra whitespace
				// 		option := strings.TrimSpace(strings.TrimPrefix(line, "Option:"))
				// 		if option != "" {
				// 			q.Options = append(q.Options, option)
				// 		}
				// 	} else if strings.HasPrefix(line, "Explanation:") {
				// 		q.Explanation = strings.TrimSpace(strings.TrimPrefix(line, "Explanation:"))
				// 	}

				// case *questionbank.FillInTheBlankQuestion:
				// 	q.BaseQuestion = baseQuestion

				// case *questionbank.TextBasedQuestion:
				// 	q.BaseQuestion = baseQuestion

				// case *questionbank.TrueFalseQuestion:
				// 	q.BaseQuestion = baseQuestion

				// }
			}

			// NOTE: Did not set specializedFields = true above as this is to be set after successful validation
			// after parsing and validating.

		}

		if !skipEntry {
			// Validate question-specific fields dynamically
			// switch q := question.(type) {
			// case *questionbank.MCQQuestion:
			// 	// Validate specialized fields after parsing
			// 	if len(q.Options) == 0 {
			// 		missingFields = append(missingFields, "Options")
			// 		skipEntry = true
			// 	}

			// 	if q.Explanation == "" {
			// 		missingFields = append(missingFields, "Explanation")
			// 		skipEntry = true
			// 	}

			// 	if !skipEntry {
			// 		specializedFields = true
			// 	}

			// case *questionbank.FillInTheBlankQuestion:
			// 	if !skipEntry {
			// 		specializedFields = true
			// 	}
			// case *questionbank.TextBasedQuestion:
			// 	if !skipEntry {
			// 		specializedFields = true
			// 	}
			// case *questionbank.TrueFalseQuestion:
			// 	if !skipEntry {
			// 		specializedFields = true
			// 	}
			// }

			// Validate presence of all common fields
			requiredCommonFields := []string{"Category", "Subcategories", "QuestionText", "Answer", "Difficulty", "Tags", "Status"}
			for _, field := range requiredCommonFields {
				if !commonFields[field] {
					missingFields = append(missingFields, field)
					skipEntry = true
				}
			}

			// Validate specialized fields
			if !specializedFields {
				missingFields = append(missingFields, "specialized fields")
				skipEntry = true
			}
		}

		// Skip adding the question if flagged
		if skipEntry {
			skippedQuestions[i] = fmt.Sprintf("Missing fields: %v", missingFields)
			continue
		}

		// // Add the question to the list if everything is valid
		// if question != nil {
		// 	questions = append(questions, question)
		// }
	}

	// Log the successful parsing of incoming questions
	log.Println("Questions Parsed Successfully.")
	log.Println(questions...)

	return questions, skippedQuestions, nil
}

// AddBulkQuestionHandler is the HTTP handler function that accepts questions and processes them.
func AddBulkQuestionHandler(c *gin.Context) {
	// Read the body of the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Ensure we close the request body when we're done
	defer c.Request.Body.Close()

	// Check if the body is empty
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
		return
	}

	// Parse the questions from the request body (assuming input is a plain text string)
	questions, skippedQuestionsAndReasons, err := parseQuestions(string(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error parsing questions: %s", err)})
		return
	}

	// Initialize counters for added and skipped questions
	totalAdded := 0
	totalSkipped := len(skippedQuestionsAndReasons)

	// Get the active PostgreSQL instance
	// db := config.GetPostgresDB()

	// Loop through the questions and insert them into the database
	// for _, question := range questions {
	// Start with the common fields of BaseQuestion
	// switch q := question.(type) {
	// case *questionbank.MCQQuestion:
	// 	if err := db.Create(&q).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error inserting MCQ question: %v", err)})
	// 		return
	// 	}
	// case *questionbank.FillInTheBlankQuestion:
	// 	if err := db.Create(&q).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add Fill-in-the-Blank question: %v", err)})
	// 		return
	// 	}
	// case *questionbank.TextBasedQuestion:
	// 	if err := db.Create(&q).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add TXT question: %v", err)})
	// 		return
	// 	}
	// case *questionbank.TrueFalseQuestion:
	// 	if err := db.Create(&q).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add True/False question: %v", err)})
	// 		return
	// 	}
	// default:
	// 	// If an unknown question type is encountered
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown question type"})
	// 	return
	// }
	// 	totalAdded++
	// }

	// Return the results as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"total_questions_added":         totalAdded,
		"total_questions_skipped":       totalSkipped,
		"skipped_questions_and_reasons": skippedQuestionsAndReasons,
		"added_questions":               questions,
	})
}

// ### Question Type: MCQ
// Category: Quantitative Aptitude
// Subcategories: Arithmetic, Basic Operations
// Difficulty: Easy
// Tags: arithmetic, simplification
// Status: active

// Question: Simplify: (10+6)×2−4÷2
// Option: 20
// Option: 30
// Option: 40
// Option: 50
// Answer: 30
// Explanation: Use BODMAS rules to simplify the expression.
// ---

// ### Question Type: FB
// Category: Grammar
// Subcategories: Fill in the Blank
// Difficulty: Medium
// Tags: grammar, sentence structure
// Status: active

// Question: The capital of France is _____.
// Answer: Paris
// ---

// ### Question Type: TXT
// Category: Logical Reasoning
// Subcategories: Puzzles
// Difficulty: Hard
// Tags: puzzles, reasoning
// Status: active

// Question: A farmer has 17 sheep, and all but 9 run away. How many are left?
// Answer: 9
// ---

// ### Question Type: TF
// Category: Science
// Subcategories: Physics
// Difficulty: Easy
// Tags: true/false, physics
// Status: active

// Question: The speed of light is constant in all mediums.
// Answer: False
// Explanation: The speed of light changes depending on the medium (e.g., slower in glass than in air).
