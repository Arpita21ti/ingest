package controllersNew

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	question_hierarchy "server/models/question_bank/question_hierarchy"
	question_type "server/models/question_bank/question_type"
	requests "server/models/requests"
	"server/models/response"
	student_psql "server/models/student_psql"
	"server/validators"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TODO: Remove the logging statements after testing and debugging.
// TODO: Remove the details from the error JSON returns after testing and debugging.

// Controller for batch loading the entire hierarchy
// GetQuestionHierarchy returns the entire question hierarchy at once,
// including domains, subdomains, niches, difficulty levels, and formats.
func GetQuestionHierarchy(c *gin.Context) {
	data := map[string]interface{}{}

	err := config.GetPostgresDBConnection().Transaction(func(tx *gorm.DB) error {
		// Fetch domains
		var domains []response.GetDomainsResponse
		if err := tx.Table("question_schema.question_domains_table").
			Select("question_domain_id", "domain_name").
			Scan(&domains).Error; err != nil {
			return fmt.Errorf("failed to fetch domains: %w", err)
		}
		data["domains"] = domains

		// Fetch subdomains
		var subDomains []response.GetSubDomainsResponse
		if err := tx.Table("question_schema.question_sub_domains_table").
			Select("question_sub_domain_id", "sub_domain_name", "question_domain_id").
			Scan(&subDomains).Error; err != nil {
			return fmt.Errorf("failed to fetch subdomains: %w", err)
		}
		data["subdomains"] = subDomains

		// Fetch niches
		var niches []response.GetNicheResponse
		if err := tx.Table("question_schema.question_niches_table").
			Select("question_niche_id", "niche_name", "question_sub_domain_id").
			Scan(&niches).Error; err != nil {
			return fmt.Errorf("failed to fetch niches: %w", err)
		}
		data["niches"] = niches

		// Fetch difficulty levels
		var difficultyLevels []response.GetDifficultyLevelsResponse
		if err := tx.Table("question_schema.question_difficulty_level_table").
			Select("question_difficulty_level_id", "difficulty_level", "question_niche_id").
			Scan(&difficultyLevels).Error; err != nil {
			return fmt.Errorf("failed to fetch difficulty levels: %w", err)
		}
		data["difficulty_levels"] = difficultyLevels

		// Fetch formats
		var formats []response.GetQuestionFormatsResponse
		if err := tx.Table("question_schema.question_formats_table").
			Select("question_format_id", "format", "question_difficulty_level_id").
			Scan(&formats).Error; err != nil {
			return fmt.Errorf("failed to fetch formats: %w", err)
		}
		data["formats"] = formats

		return nil // Commit the transaction
	})

	if err != nil {
		// Return the error if the transaction fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send all data in a single response
	c.JSON(http.StatusOK, data)
}

// GetDomains returns all question domains
func GetDomains(c *gin.Context) {
	var domains []response.GetDomainsResponse

	// Fetch only domainId and domainName
	if err := config.GetPostgresTable(&question_hierarchy.QuestionDomainsTable{}).
		Select("question_domain_id", "domain_name"). // Select only the required columns for added security.
		Scan(&domains).                              // Map the result to the DTO.
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains", "details": err.Error()})
		return
	}

	// Check if no domains were found
	if len(domains) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No domains found"})
		return
	}

	// Log the number of records fetched
	log.Printf("Fetched %d domains", len(domains))

	// Optionally, you can also log the size of the data in bytes
	dataSize := len(domains) * int(unsafe.Sizeof(domains[0]))
	log.Printf("Data size (approximate): %d bytes", dataSize)

	// Send the response
	c.JSON(http.StatusOK, domains)
}

// GetSubDomains returns all question sub-domains
func GetSubDomains(c *gin.Context) {

	// Validate the presence of domainID in the URL query
	domainId := c.Param("domainID")
	if domainId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: domainId"})
		return
	}

	var subDomains []response.GetSubDomainsResponse

	// Fetch only domainId and domainName
	if err := config.GetPostgresTable(&question_hierarchy.QuestionSubDomainsTable{}).
		Select("question_sub_domain_id", "sub_domain_name", "question_domain_id"). // Select only the required columns for added security.
		Where("question_domain_id = ?", domainId).                                 // Filter the sub-domains based on the domain ID.
		Scan(&subDomains).                                                         // Map the result to the DTO.
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub domains", "details": err.Error()})
		return
	}

	// Check if no subDomains were found
	if len(subDomains) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No subDomains found"})
		return
	}

	// Log the number of records fetched
	log.Printf("Fetched %d subDomains", len(subDomains))

	// Optionally, you can also log the size of the data in bytes
	dataSize := len(subDomains) * int(unsafe.Sizeof(subDomains[0]))
	log.Printf("Data size (approximate): %d bytes", dataSize)

	// Send the response
	c.JSON(http.StatusOK, subDomains)
}

// GetNiches returns all question niches
func GetNiches(c *gin.Context) {

	// Validate the presence of subDomainsID in the URL query
	subDomainsID := c.Param("subDomainsID")
	if subDomainsID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: subDomainsID"})
		return
	}

	var niches []response.GetNicheResponse

	// Fetch only domainId and domainName
	if err := config.GetPostgresTable(&question_hierarchy.QuestionNicheTable{}).
		Select("question_niche_id", "niche_name").         // Select only the required columns for added security.
		Where("question_sub_domain_id = ?", subDomainsID). // Filter the sub-domains based on the domain ID.
		Scan(&niches).                                     // Map the result to the DTO.
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub domains", "details": err.Error()})
		return
	}

	// Check if no niches were found
	if len(niches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No niches found"})
		return
	}

	// Log the number of records fetched
	log.Printf("Fetched %d niches", len(niches))

	// Optionally, you can also log the size of the data in bytes
	dataSize := len(niches) * int(unsafe.Sizeof(niches[0]))
	log.Printf("Data size (approximate): %d bytes", dataSize)

	c.JSON(http.StatusOK, niches)
}

// GetDifficultyLevels returns all question difficulty levels
func GetDifficultyLevels(c *gin.Context) {
	// Validate the presence of nicheID in the URL query
	nicheID := c.Param("nicheID")
	if nicheID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: nicheID"})
		return
	}

	var difficultyLevels []response.GetDifficultyLevelsResponse

	// Fetch only domainId and domainName
	if err := config.GetPostgresTable(&question_hierarchy.QuestionDifficultyLevelTable{}).
		Select("question_difficulty_level_id", "difficulty_level"). // Select only the required columns for added security.
		Where("question_niche_id = ?", nicheID).                    // Filter the sub-domains based on the domain ID.
		Scan(&difficultyLevels).                                    // Map the result to the DTO.
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Difficulty Levels", "details": err.Error()})
		return
	}

	// Check if no difficultyLevels were found
	if len(difficultyLevels) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No difficulty levels found"})
		return
	}

	// Log the number of records fetched
	log.Printf("Fetched %d difficulty levels", len(difficultyLevels))

	// Optionally, you can also log the size of the data in bytes
	dataSize := len(difficultyLevels) * int(unsafe.Sizeof(difficultyLevels[0]))
	log.Printf("Data size (approximate): %d bytes", dataSize)

	c.JSON(http.StatusOK, difficultyLevels)
}

// GetFormats returns all question formats
func GetFormats(c *gin.Context) {
	// Validate the presence of difficultyLevelID in the URL query
	difficultyLevelID := c.Param("difficultyLevelID")
	if difficultyLevelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: difficultyLevelID"})
		return
	}

	var formats []response.GetQuestionFormatsResponse

	// Fetch only domainId and domainName
	if err := config.GetPostgresTable(&question_hierarchy.QuestionFormatTable{}).
		Select("question_format_id", "format").                       // Select only the required columns for added security.
		Where("question_difficulty_level_id = ?", difficultyLevelID). // Filter the sub-domains based on the domain ID.
		Scan(&formats).                                               // Map the result to the DTO.
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Difficulty Levels", "details": err.Error()})
		return
	}

	// Check if no formats were found
	if len(formats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No formats found"})
		return
	}

	// Log the number of records fetched
	log.Printf("Fetched %d formats", len(formats))

	// Optionally, you can also log the size of the data in bytes
	dataSize := len(formats) * int(unsafe.Sizeof(formats[0]))
	log.Printf("Data size (approximate): %d bytes", dataSize)

	c.JSON(http.StatusOK, formats)
}

// Helper function to validate practiceSession Record input
func validateStudentPracticeSessionRecordTableInput(input requests.GetQuestionsRequest) error {
	validate := validator.New()
	validators.RegisterValidatorsPracticeSession(validate)
	return validate.Struct(input)
}

// Generic helper function to fetch questions of any type
func fetchQuestionsOfType[T any](formatId, lastAttemptedQuestionID uint32, requiredQuestionCount int) ([]T, error) {
	var questions []T
	if err := config.GetPostgresDBConnection().
		Where("question_format_id = ? AND id > ?", formatId, lastAttemptedQuestionID).
		Limit(requiredQuestionCount).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %v", err)
	}

	// If questions are less than required, fetch remaining questions
	if len(questions) < requiredQuestionCount {
		var remainingQuestions []T
		remainingCount := requiredQuestionCount - len(questions)
		if err := config.GetPostgresDBConnection().
			Where("question_format_id = ? AND id <= ?", formatId, 1).
			Limit(remainingCount).
			Find(&remainingQuestions).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch additional questions: %v", err)
		}
		questions = append(questions, remainingQuestions...)
	}

	return questions, nil
}

// Helper function to fetch questions based on format
func fetchQuestionsByFormat(questionFormat string, formatId, lastAttemptedQuestionID uint32, requiredQuestionCount int) (interface{}, error) {
	var questions interface{}
	var err error

	switch questionFormat {
	case "MCQ":
		questions, err = fetchQuestionsOfType[question_type.MCQQuestion](formatId, lastAttemptedQuestionID, requiredQuestionCount)
	case "TF":
		questions, err = fetchQuestionsOfType[question_type.TrueFalseQuestion](formatId, lastAttemptedQuestionID, requiredQuestionCount)
	case "FIB":
		questions, err = fetchQuestionsOfType[question_type.FillInTheBlankQuestion](formatId, lastAttemptedQuestionID, requiredQuestionCount)
	case "TXT":
		questions, err = fetchQuestionsOfType[question_type.TextBasedQuestion](formatId, lastAttemptedQuestionID, requiredQuestionCount)
	default:
		return nil, fmt.Errorf("invalid question format")
	}

	return questions, err
}

// GetQuestions returns questions in a paginated way for practice sessions
func GetQuestions(c *gin.Context) {

	// Parse the incoming request body
	var request requests.GetQuestionsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate the enrollment number
	if err := validateStudentPracticeSessionRecordTableInput(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Extract values from the request
	requiredQuestionCount := request.QuestionCount
	formatId := request.QuestionFormatID
	lastAttemptedQuestionID := request.LastAttemptedQuestionID

	// Validate the number of questions to attempt limit
	if requiredQuestionCount != 1 && requiredQuestionCount != 10 && requiredQuestionCount != 30 && requiredQuestionCount != 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QuestionCount. Must be 10, 30, or 60"})
		return
	}

	// Fetch random questions from the appropriate table based on format
	questions, err := fetchQuestionsByFormat(request.QuestionFormat, formatId, lastAttemptedQuestionID, requiredQuestionCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// No need to check for an active session because if user terminates the session halfway the front-end
	// is required to forcefully end the session using the "/end-forcefully" route.

	// Store the practice session record
	practiceSessionRecord := student_psql.StudentPracticeSessionRecordTable{
		DomainID:           request.QuestionDomainID,
		SubDomainID:        request.QuestionSubDomainID,
		DifficultyLevelID:  request.QuestionDifficultyLevelID,
		QuestionsAttempted: -1,                              // This will be updated after the session is completed
		QuestionsCorrect:   -1,                              // This will be updated after the session is completed
		ScoreEarned:        -1,                              // This will be updated after the session is completed
		StartTime:          time.Now().Add(5 * time.Second), // Adding 5 seconds to the current time to simulate the start time considering the time taken to fetch and transfer the questions
		EndTime:            time.Time{},                     // Default value indicating the end time is not set yet
	}

	if err := config.GetPostgresDBConnection().Create(&practiceSessionRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store practice session record", "details": err.Error()})
		return
	}

	// Insert the record in the practice session lookup table
	practiceSessionLookupRecord := student_psql.StudentPracticeSessionLookupTable{
		EnrollmentNo:      request.EnrollmentNo,
		PracticeSessionID: practiceSessionRecord.PracticeSessionID,
	}

	if err := config.GetPostgresDBConnection().Create(&practiceSessionLookupRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store practice session record in lookup", "details": err.Error()})
		return

	}

	response := response.GetQuestionsResponse{
		Questions:         questions,
		PracticeSessionID: practiceSessionRecord.PracticeSessionID,
		Message:           "Practice session started successfully",
	}

	// Return the questions and practiceSessionID in the response
	c.JSON(
		http.StatusOK,
		response,
	)
}

// // GetQuestions returns questions in a paginated way for practice sessions
// func GetQuestions(c *gin.Context) {

// 	// Parse the incoming request body
// 	var request requests.GetQuestionsRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
// 		return
// 	}

// 	// Validate the enrollment number
// 	if err := validateStudentPracticeSessionRecordTableInput(request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
// 		return
// 	}

// 	// Extract values from the request
// 	requiredQuestionCount := request.QuestionCount
// 	formatId := request.QuestionFormatID
// 	lastAttemptedQuestionID := request.LastAttemptedQuestionID

// 	// Validate the number of questions to attempt limit
// 	if requiredQuestionCount != 1 && requiredQuestionCount != 10 && requiredQuestionCount != 30 && requiredQuestionCount != 60 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QuestionCount. Must be 10, 30, or 60"})
// 		return
// 	}

// 	// No need to check for an active session because if user terminates the session halfway the front-end
// 	// is required to forcefully end the session using the "/end-forcefully" route.

// 	// Fetch random questions from the appropriate table based on format
// 	var questions interface{}
// 	switch request.QuestionFormat {
// 	case "MCQ":
// 		var mcqQuestions []question_type.MCQQuestion
// 		if err := config.GetPostgresDBConnection().
// 			Where("question_format_id = ? AND id > ?", formatId, lastAttemptedQuestionID).
// 			Limit(requiredQuestionCount).
// 			Find(&mcqQuestions).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch MCQ questions", "details": err.Error()})
// 			return
// 		}
// 		if len(mcqQuestions) != requiredQuestionCount {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the required number of MCQ questions"})
// 			return
// 		}
// 		questions = mcqQuestions

// 	case "TF":
// 		var tfQuestions []question_type.TrueFalseQuestion
// 		if err := config.GetPostgresDBConnection().
// 			Where("question_format_id = ? AND id > ?", formatId, lastAttemptedQuestionID).
// 			Limit(requiredQuestionCount).
// 			Find(&tfQuestions).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch True/False questions", "details": err.Error()})
// 			return
// 		}
// 		if len(tfQuestions) != requiredQuestionCount {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the required number of True/False questions"})
// 			return
// 		}
// 		questions = tfQuestions

// 	case "FIB":
// 		var fibQuestions []question_type.FillInTheBlankQuestion
// 		if err := config.GetPostgresDBConnection().
// 			Where("question_format_id = ? AND id > ?", formatId, lastAttemptedQuestionID).
// 			Limit(requiredQuestionCount).
// 			Find(&fibQuestions).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Fill-in-the-Blank questions", "details": err.Error()})
// 			return
// 		}
// 		if len(fibQuestions) != requiredQuestionCount {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the required number of Fill-in-the-Blank questions"})
// 			return
// 		}
// 		questions = fibQuestions

// 	case "TXT":
// 		var textQuestions []question_type.TextBasedQuestion
// 		if err := config.GetPostgresDBConnection().
// 			Where("question_format_id = ? AND id > ?", formatId, lastAttemptedQuestionID).
// 			Limit(requiredQuestionCount).
// 			Find(&textQuestions).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Text-based questions", "details": err.Error()})
// 			return
// 		}
// 		if len(textQuestions) != requiredQuestionCount {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the required number of Text-based questions"})
// 			return
// 		}
// 		questions = textQuestions

// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question format"})
// 		return
// 	}

// 	// Store the practice session record
// 	practiceSessionRecord := student_psql.StudentPracticeSessionRecordTable{
// 		DomainID:           request.QuestionDomainID,
// 		SubDomainID:        request.QuestionSubDomainID,
// 		DifficultyLevelID:  request.QuestionDifficultyLevelID,
// 		QuestionsAttempted: -1,                              // This will be updated after the session is completed
// 		QuestionsCorrect:   -1,                              // This will be updated after the session is completed
// 		ScoreEarned:        -1,                              // This will be updated after the session is completed
// 		StartTime:          time.Now().Add(5 * time.Second), // Adding 5 seconds to the current time to simulate the start time considering the time taken to fetch and transfer the questions
// 		EndTime:            time.Time{},                     // Default value indicating the end time is not set yet
// 	}

// 	if err := config.GetPostgresDBConnection().Create(&practiceSessionRecord).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store practice session record", "details": err.Error()})
// 		return
// 	}

// 	// Insert the record in the practice session lookup table
// 	practiceSessionLookupRecord := student_psql.StudentPracticeSessionLookupTable{
// 		EnrollmentNo:      request.EnrollmentNo,
// 		PracticeSessionID: practiceSessionRecord.PracticeSessionID,
// 	}

// 	if err := config.GetPostgresDBConnection().Create(&practiceSessionLookupRecord).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store practice session record in lookup", "details": err.Error()})
// 		return

// 	}

// 	response := response.GetQuestionsResponse{
// 		Questions:         questions,
// 		PracticeSessionID: practiceSessionRecord.PracticeSessionID,
// 		Message:           "Practice session started successfully",
// 	}

// 	// Return the questions and practiceSessionID in the response
// 	c.JSON(
// 		http.StatusOK,
// 		response,
// 	)
// }
