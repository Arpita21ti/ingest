package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"server/config"
	models "server/models/student"
	"server/utils"
	"server/validators"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to validate password reset input
func validatePasswordResetInput(input models.StudentForgotPasswordRequest) error {
	validate := validator.New()
	validators.RegisterValidatorsPasswordReset(validate) // Register custom validators here
	return validate.Struct(input)
}

// Helper function to get the collection for StudentPasswordResetRequests
func getStudentPasswordResetRequestsCollection() *mongo.Collection {
	return config.GetCollectionMongo("StudentPasswordResetRequests")
}

// updateResetToken function updates or creates a reset token and its expiry for a student
func updateResetToken(enrollmentNo, token string, expiry time.Time) error {
	collection := getStudentPasswordResetRequestsCollection()

	// Create a PasswordResetRequest model instance
	resetRequest := models.StudentPasswordResetRequestsStorage{
		EnrollmentNo:     enrollmentNo,
		ResetToken:       token,
		ResetTokenExpiry: expiry,
		UpdatedAt:        time.Now(),
	}

	// Define the filter to locate the student by enrollment number
	filter := bson.M{"enrollmentNo": enrollmentNo}

	// Define the update operation using the model
	update := bson.M{"$set": resetRequest}

	// Set the options to upsert (create if not exists)
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

// verifyResetToken verifies the token and retrieves the corresponding user.
func verifyResetToken(enrollmentNo, token string) (bool, error) {
	collection := getStudentPasswordResetRequestsCollection()

	// Find the user by the reset token
	filter := bson.M{
		"enrollmentNo":       enrollmentNo, // Ensures the same user has sent the request
		"reset_token":        token,
		"reset_token_expiry": bson.M{"$gt": time.Now()}, // Ensure the token is not expired
	}

	var resetRequest models.StudentPasswordResetRequestsStorage
	err := collection.FindOne(context.TODO(), filter).Decode(&resetRequest)
	if err == mongo.ErrNoDocuments {
		return false, errors.New("invalid or expired token")
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// updatePasswordAndClearToken updates the password for a user and clears the reset token
func updatePasswordAndClearToken(enrollmentNo, hashedPassword string) error {
	// Get collections
	resetTokenStoreCollection := getStudentPasswordResetRequestsCollection()
	studentLoginCredentialsCollection := getStudentLoginCredentialsCollection()

	// Start a session for transactional operation
	session, err := config.ClientMongo.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	// Use the session for a transaction
	err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {
		// Begin transaction
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Step 1: Update the password in the login credentials collection
		loginFilter := bson.M{"enrollmentNo": enrollmentNo}
		loginUpdate := bson.M{
			"$set": bson.M{
				"password": hashedPassword,
			},
		}

		if _, err := studentLoginCredentialsCollection.UpdateOne(sc, loginFilter, loginUpdate); err != nil {
			session.AbortTransaction(sc) // Abort the transaction if any error occurs
			return err
		}

		// Step 2: Remove the reset token entry from the reset token store collection
		resetTokenFilter := bson.M{"enrollmentNo": enrollmentNo}
		if _, err := resetTokenStoreCollection.DeleteOne(sc, resetTokenFilter); err != nil {
			session.AbortTransaction(sc) // Abort the transaction if any error occurs
			return err
		}

		// Commit the transaction
		if err := session.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	})

	return err
}

func ResetPasswordHandler(c *gin.Context) {
	// Define the request payload for the new password
	var req models.StudentResetPasswordRequest

	// Parse and validate the new password from the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	// Verify the token follows the Bearer format
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
		return
	}

	// Extract the token value
	token := authHeader[len(bearerPrefix):]

	// Verify token and fetch user
	user, err := verifyResetToken(req.EnrollmentNo, token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying token", "details": err.Error()})
		return
	} else if !user {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	// Update password and clear the reset token
	if err := updatePasswordAndClearToken(req.EnrollmentNo, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func ForgotPasswordHandler(c *gin.Context) {

	var userInput models.StudentForgotPasswordRequest

	// Parse and validate input
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate user input with custom validators
	if err := validatePasswordResetInput(userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Retrieve user from the database
	user, err := fetchUserByEnrollmentNo(userInput.EnrollmentNo)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) //Didn't return Invalid Enrollment No. for more security.
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user", "details": err.Error()})
		return
	}

	// Verify the email
	if userInput.Email != user.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) //For security, did't specify the reason.
		return
	}

	// Verify the phone number
	if userInput.Phone != user.Phone {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) //For security, did't specify the reason.
		return
	}

	// Generate JWT token for student user
	commonTokenString, err := utils.GenerateToken("student_user", "common", 24)
	if err != nil {
		fmt.Println("Error generating common user token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	// Save token and expiry to the database
	expiry := time.Now().Add(15 * time.Minute)
	if err := updateResetToken(user.Email, commonTokenString, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token", "details": err.Error()})
		return
	}

	// Send email with the reset link
	if err := utils.SendResetEmail(user.Email, commonTokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
}
