package controllersNew

import (
	"errors"
	"fmt"
	"net/http"
	"server/config"

	requests "server/models/requests"
	student_psql "server/models/student_psql"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NOTE: Session Start is taken care of by the GetQuestions handler in the question_controller.go

// SubmitPracticeSessionHandler submits the results of a practice session
func SubmitPracticeSessionHandler(c *gin.Context) {

	var request requests.SuccessfullyEndPracticeSessionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Use the transaction method
	err := config.GetPostgresDBConnection().Transaction(func(tx *gorm.DB) error {

		// Fetch and update the lookup table record
		var practiceSessionLookupRecord student_psql.StudentPracticeSessionLookupTable
		if err := tx.Where("practice_session_id = ? AND status = ?", request.PracticeSessionID, "Active").
			First(&practiceSessionLookupRecord).Error; err != nil {
			return fmt.Errorf("no active practice session found: %w", err)
		}

		practiceSessionLookupRecord.Status = "Submitted"

		if err := tx.Save(&practiceSessionLookupRecord).Error; err != nil {
			return fmt.Errorf("failed to update practice session status: %w", err)
		}

		// Fetch and update the practice session record
		var practiceSessionRecord student_psql.StudentPracticeSessionRecordTable
		if err := tx.Where("practice_session_id = ?", request.PracticeSessionID).
			First(&practiceSessionRecord).Error; err != nil {
			return fmt.Errorf("practice session not found: %w", err)
		}

		practiceSessionRecord.QuestionsAttempted = request.QuestionsAttempted
		practiceSessionRecord.QuestionsCorrect = request.QuestionsCorrect
		practiceSessionRecord.ScoreEarned = request.ScoreEarnedPercentage
		practiceSessionRecord.EndTime = time.Now()
		practiceSessionRecord.Feedbacks = request.Feedbacks

		if err := tx.Save(&practiceSessionRecord).Error; err != nil {
			return fmt.Errorf("failed to submit practice session: %w", err)
		}

		return nil
	})

	// Handle the result of the transaction
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Practice session submitted successfully"})
}

// ForcefullyEndPracticeSessionHandler forcefully ends a practice session
func ForcefullyEndPracticeSessionHandler(c *gin.Context) {
	var request requests.ForcefullyEndPracticeSessionRequest

	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Use the transaction method
	err := config.GetPostgresDBConnection().Transaction(func(tx *gorm.DB) error {
		// Find the practice session record
		var practiceSessionRecord student_psql.StudentPracticeSessionLookupTable
		if err := tx.Where("practice_session_id = ? AND status = ?", request.PracticeSessionID, "Active").
			First(&practiceSessionRecord).Error; err != nil {
			return fmt.Errorf("no active practice session found: %w", err)
		}

		// Update the status of the practice session
		practiceSessionRecord.Status = "Force End"
		if err := tx.Save(&practiceSessionRecord).Error; err != nil {
			return fmt.Errorf("failed to update practice session status: %w", err)
		}

		// Delete the corresponding record in the StudentPracticeSessionRecordTable
		if err := tx.Where("practice_session_id = ?", practiceSessionRecord.PracticeSessionID).
			Delete(&student_psql.StudentPracticeSessionRecordTable{}).Error; err != nil {
			return fmt.Errorf("failed to delete practice session record: %w", err)
		}

		return nil
	})

	// Handle the result of the transaction
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Practice session forcefully ended successfully"})
}
