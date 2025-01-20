package controllers

import (
	"context"
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
	"golang.org/x/crypto/bcrypt"
)

// Helper function to get the collection for StudentProfile
func getStudentProfileCollection() *mongo.Collection {
	return config.GetCollectionMongo("StudentProfile")
}

// Helper function to get the collection for StudentLoginCredentials
func getStudentLoginCredentialsCollection() *mongo.Collection {
	return config.GetCollectionMongo("StudentLoginCredentials")
}

// Helper function to get the collection for StudentCollegeAcademics
func getStudentAcademicsCollection() *mongo.Collection {
	return config.GetCollectionMongo("StudentCollegeAcademics")
}

// Helper function to get the collection for StudentCertificationsAndResume
func getStudentCertificationsCollection() *mongo.Collection {
	return config.GetCollectionMongo("StudentCertificationsAndResume")
}

// Helper function to validate login input
func validateLoginInput(input models.StudentLogin) error {
	validate := validator.New()
	validators.RegisterValidatorsLogIn(validate) // Register custom validators here
	return validate.Struct(input)
}

// Helper function to validate signup input
func validateSignUpInput(input models.StudentSignUp) error {
	validate := validator.New()
	validators.RegisterValidatorsSignUp(validate)
	return validate.Struct(input)
}

// Helper function to fetch user by enrollment number
func fetchUserByEnrollmentNo(enrollmentNo string) (*models.StudentLogin, error) {
	var user models.StudentLogin
	collection := getStudentLoginCredentialsCollection()
	err := collection.FindOne(context.Background(), bson.M{"enrollmentNo": enrollmentNo}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Helper function to check if a record exists based on the EnrollmentNo
func enrollmentNoRecordExists(ctx context.Context, enrollmentNo string) (bool, error) {
	// Example: Check the StudentLoginCredentials collection for existing record
	loginCollection := getStudentLoginCredentialsCollection()
	var result models.StudentLogin
	err := loginCollection.FindOne(ctx, bson.M{"enrollmentNo": enrollmentNo}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No record found
			return false, nil
		}
		// Some other error occurred
		return false, err
	}
	// Record exists
	return true, nil
}

// Helper function for MongoDB insertion
func insertIntoCollection(ctx context.Context, collection *mongo.Collection, document interface{}, entityName string) error {
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		// Log the error (replace with your logging library)
		fmt.Printf("Failed to insert %s: %v\n", entityName, err)
		return fmt.Errorf("failed to create %s", entityName)
	}
	return nil
}

// StudentLoginHandler handler
func StudentLoginHandler(c *gin.Context) {

	var userInput models.StudentLogin

	// Parse and validate input
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Validate user input with custom validators
	if err := validateLoginInput(userInput); err != nil {
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

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) //Didn't return Invalid Enrollment No. for more security.
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

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{"Authorization": "Bearer " + commonTokenString})
}

// StudentSignupHandler handler
func StudentSignupHandler(c *gin.Context) {
	// Define a timeout context for all database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userInput models.StudentSignUp

	// Parse input
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate the input
	if err := validateSignUpInput(userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Check if student record with the same EnrollmentNo already exists
	if exists, err := enrollmentNoRecordExists(ctx, userInput.EnrollmentNo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking record existence", "details": err.Error()})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Student with this EnrollmentNo already exists"})
		return
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	userInput.Password = string(hashedPassword)

	// Prepared data for insertion to individual collections

	// Create the student login credentials
	studentLogin := models.StudentLogin{
		EnrollmentNo: userInput.EnrollmentNo,
		Email:        userInput.Email,
		Password:     userInput.Password,
		Phone:        userInput.Phone,
	}

	// Create the student profile
	studentProfile := models.StudentProfile{
		EnrollmentNo:       userInput.EnrollmentNo,
		Name:               userInput.Name,
		Branch:             userInput.Branch,
		YearOfEnrollment:   userInput.YearOfEnrollment,
		ClassTenPercent:    userInput.ClassTenPercent,
		ClassTwelvePercent: userInput.ClassTwelvePercent,
	}

	// Create the student academic details
	studentAcademics := models.StudentCollegeAcademics{
		EnrollmentNo:    userInput.EnrollmentNo,
		PreviousSemCGPA: userInput.PreviousSemCGPA,
		PreviousSemSGPA: userInput.PreviousSemSGPA,
	}

	// Create the student certifications and resume
	studentCertifications := models.StudentCertificationsAndResume{
		EnrollmentNo:   userInput.EnrollmentNo,
		Certifications: userInput.Certifications,
		Resume:         userInput.Resume,
	}

	// Saveing data to MongoDB (use a helper function for reusability)

	// Saved student login credentials in the "StudentLoginCredentials" collection
	if err := insertIntoCollection(ctx, getStudentLoginCredentialsCollection(), studentLogin, "login credentials"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Saved student profile in the "StudentProfile" collection
	if err := insertIntoCollection(ctx, getStudentProfileCollection(), studentProfile, "student profile"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save student academics in the "StudentCollegeAcademics" collection
	if err := insertIntoCollection(ctx, getStudentAcademicsCollection(), studentAcademics, "student academics"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save student certifications in the "StudentCertificationsAndResume" collection
	if err := insertIntoCollection(ctx, getStudentCertificationsCollection(), studentCertifications, "student certifications and resume"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token for student user
	commonTokenString, err := utils.GenerateToken("student_user", "common", 24)
	if err != nil {
		fmt.Println("Error generating common user token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	// Respond with success and JWT
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Student Profile created successfully",
		"Authorization": "Bearer " + commonTokenString,
	})
}
