package controllersNew

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"server/config"
	"strconv"

	requests "server/models/requests"
	models "server/models/student_psql"
	"server/utils"
	"server/validators"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper function to validate login input
func validateLoginInput(input requests.StudentLoginRequest) error {
	validate := validator.New()
	validators.RegisterValidatorsLogIn(validate) // Register custom validators here
	return validate.Struct(input)
}

// Helper function to fetch user by enrollment number
func fetchUserByEnrollmentNo(enrollmentNo string) (*models.EnrollmentMasterLookupTable, error) {
	var user models.EnrollmentMasterLookupTable
	if err := config.GetPostgresDBConnection().Where("enrollment_no = ?", enrollmentNo).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// StudentLoginHandler handler for student login
func StudentLoginHandler(c *gin.Context) {
	var userInput requests.StudentLoginRequest

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

	// Retrieve user from the database based on enrollment number
	user, err := fetchUserByEnrollmentNo(userInput.EnrollmentNo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Didn't return Invalid Enrollment No for more security.
		return
	}

	// Fetch login details using the LogInDetailsID
	var loginDetails models.StudentLogInDetailsTable
	if err := config.GetPostgresDBConnection().Where("id = ?", user.LogInDetailsID).First(&loginDetails).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Enrollment"}) // Didn't return Invalid LogInDetailsID for more security.
		return
	}

	// Verify the email
	if userInput.Email != loginDetails.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email"}) // For security, didn't specify the reason.
		return
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(loginDetails.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"}) // Didn't return Invalid Password for more security.
		return
	}

	// Verify the phone number
	if userInput.Phone != loginDetails.Phone {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Phone"}) // For security, didn't specify the reason.
		return
	}

	// Generate JWT token for student user
	commonTokenString, err := utils.GenerateToken("student_user", "common", 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	// Respond with the token and success message
	c.JSON(
		http.StatusOK, gin.H{
			"message":       "Login successfully",
			"Authorization": "Bearer " + commonTokenString,
		},
	)
}

// Helper function to validate signup input
func validateSignupInput(input requests.StudentSignUpRequest) error {
	validate := validator.New()
	validators.RegisterValidatorsSignUp(validate) // Register custom validators here
	return validate.Struct(input)
}

// StudentSignupHandler handler for student signup
func StudentSignupHandler(c *gin.Context) {
	var userInput requests.StudentSignUpRequest

	// Parse and validate input from the multipart form
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate user input with custom validators
	if err := validateSignupInput(userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Check if the enrollment number already exists
	if _, err := fetchUserByEnrollmentNo(userInput.EnrollmentNo); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Enrollment number already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	// Assumed file validations are done in the frontend.
	// Process the files (marksheets, photograph, resume, certifications, etc.)
	classTenFile, _ := c.FormFile("classTenMarksheet")
	classTwelveFile, _ := c.FormFile("classTwelveMarksheet")
	photographFile, _ := c.FormFile("photograph")
	resumeFile, _ := c.FormFile("resume")
	certificationFiles := c.Request.MultipartForm.File["certifications"]

	var (
		classTenMarksheetURL    string
		classTwelveMarksheetURL string
		photographURL           string
		resumeURL               string
	)

	// Upload the Class Ten Marksheet
	if classTenFile != nil {
		classTenMarksheetURL, err = uploadFileToCloud(classTenFile, userInput.EnrollmentNo+"_classTenMarksheet")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload class ten marksheet", "details": err.Error()})
			return
		}
	}

	// Upload the Class Twelve Marksheet
	if classTwelveFile != nil {
		classTwelveMarksheetURL, err = uploadFileToCloud(classTwelveFile, userInput.EnrollmentNo+"_classTwelveMarksheet")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload class twelve marksheet", "details": err.Error()})
			return
		}
	}

	// Upload the Photograph
	if photographFile != nil {
		photographURL, err = uploadFileToCloud(photographFile, userInput.EnrollmentNo+"_photograph")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photograph", "details": err.Error()})
			return
		}
	}

	// Upload the Resume
	if resumeFile != nil {
		resumeURL, err = uploadFileToCloud(resumeFile, userInput.EnrollmentNo+"_resume")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload resume", "details": err.Error()})
			return
		}
	}

	// Initialize the database connection
	db := config.GetPostgresDBConnection()

	// Start a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		// Create login details
		loginDetails := models.StudentLogInDetailsTable{
			Email:    userInput.Email,
			Password: string(hashedPassword),
			Phone:    userInput.Phone,
		}

		if err := tx.Create(&loginDetails).Error; err != nil {
			return fmt.Errorf("failed to create login details: %w", err)
		}

		// Create uploaded document details for class ten marksheet
		var classTenDocumentID, classTwelveDocumentID uint
		if classTenFile != nil {
			studentDoc := models.StudentDocumentTable{
				StoredIn:     "AWSS3",
				DocumentType: "classTenMarksheet",
				URL:          classTenMarksheetURL,
			}

			if err := tx.Create(&studentDoc).Error; err != nil {
				return fmt.Errorf("failed to create document entry for class ten: %w", err)
			}
			classTenDocumentID = uint(studentDoc.DocumentID)
		}

		// Create uploaded document details for class twelve marksheet
		if classTwelveFile != nil {
			studentDoc := models.StudentDocumentTable{
				StoredIn:     "AWSS3",
				DocumentType: "classTwelveMarksheet",
				URL:          classTwelveMarksheetURL,
			}

			if err := tx.Create(&studentDoc).Error; err != nil {
				return fmt.Errorf("failed to create document entry for class twelve: %w", err)
			}
			classTwelveDocumentID = uint(studentDoc.DocumentID)
		}

		// Create academic details with the uploaded file URLs
		academicDetails := models.StudentAcademicDetailsTable{
			Branch:                 userInput.Branch,
			YearOfEnrollment:       userInput.YearOfEnrollment,
			CGPA:                   userInput.CGPA,
			PreviousSemSGPA:        userInput.PreviousSemSGPA,
			SchoolForClassTen:      userInput.SchoolForClassTen,
			ClassTenPercentage:     userInput.ClassTenPercentage,
			ClassTenMarksheetID:    uint32(classTenDocumentID),
			SchoolForClassTwelve:   userInput.SchoolForClassTwelve,
			ClassTwelvePercentage:  userInput.ClassTwelvePercentage,
			ClassTwelveMarksheetID: uint32(classTwelveDocumentID),
		}

		if err := tx.Create(&academicDetails).Error; err != nil {
			return fmt.Errorf("failed to create academic details: %w", err)
		}

		// Create family details
		familyDetails := models.StudentFamilyDetailsTable{
			FatherName:          userInput.FatherName,
			FatherQualification: userInput.FatherQualification,
			FatherProfession:    userInput.FatherProfession,
			MotherName:          userInput.MotherName,
			MotherQualification: userInput.MotherQualification,
			MotherProfession:    userInput.MotherProfession,
			NoOfSiblings:        userInput.NoOfSiblings,
			TotalFamilyIncome:   userInput.TotalFamilyIncome,
		}

		if err := tx.Create(&familyDetails).Error; err != nil {
			return fmt.Errorf("failed to create family details: %w", err)
		}

		var photographDocumentID, resumeDocumentID uint
		// Create uploaded document details for photograph
		if photographFile != nil {
			studentDoc := models.StudentDocumentTable{
				StoredIn:     "AWSS3",
				DocumentType: "photograph",
				URL:          photographURL,
			}

			if err := tx.Create(&studentDoc).Error; err != nil {
				return fmt.Errorf("failed to create document entry for photograph: %w", err)
			}
			photographDocumentID = uint(studentDoc.DocumentID)
		}

		// Create uploaded document details for resume
		if resumeFile != nil {
			studentDoc := models.StudentDocumentTable{
				StoredIn:     "AWSS3",
				DocumentType: "resume",
				URL:          resumeURL,
			}

			if err := tx.Create(&studentDoc).Error; err != nil {
				return fmt.Errorf("failed to create document entry for resume: %w", err)
			}
			resumeDocumentID = uint(studentDoc.DocumentID)
		}

		// Create profile details
		profileDetails := models.StudentProfileDetailsTable{
			Name:         userInput.Name,
			Gender:       userInput.Gender,
			Category:     userInput.Category,
			PhotographID: uint32(photographDocumentID),
			ResumeID:     uint32(resumeDocumentID),
		}

		if err := tx.Create(&profileDetails).Error; err != nil {
			return fmt.Errorf("failed to create profile details: %w", err)
		}

		// Create scholarship details
		scholarshipDetails := models.StudentScholarshipDetailsTable{
			ScholarshipName: userInput.ScholarshipName,
			ProvidedBy:      userInput.ProvidedBy,
			AmountReceived:  userInput.AmountReceived,
		}

		if err := tx.Create(&scholarshipDetails).Error; err != nil {
			return fmt.Errorf("failed to create scholarship details: %w", err)
		}

		// Create uploaded document details for certifications
		for i, file := range certificationFiles {
			var certificationDocumentID uint
			// Upload certification file to cloud and get the URL
			certificationURL, err := uploadFileToCloud(file, userInput.EnrollmentNo+"_certificate_"+strconv.Itoa(i))
			if err != nil {
				return fmt.Errorf("failed to upload certification %s: %w", file.Filename, err)
			}

			studentDoc := models.StudentDocumentTable{
				StoredIn:     "AWSS3",
				DocumentType: userInput.EnrollmentNo + "certificate",
				URL:          certificationURL,
			}

			if err := tx.Create(&studentDoc).Error; err != nil {
				return fmt.Errorf("failed to create document entry for certificate%d: %w", i, err)
			}
			certificationDocumentID = uint(studentDoc.DocumentID)

			// Fetch the certification details (assuming they come from the request)
			certificationName := userInput.CertificationNames[i] // Example: Get certification name from the request
			issuingAuthority := userInput.IssuingAuthority[i]    // Example: Get issuing authority from the request
			issuingDate := userInput.IssuingDate[i]              // Example: Get issuing date from the request

			// Check if issuingDate is in the format YYYY-MM
			var issuingYear, issuingMonth int
			if _, err := fmt.Sscanf(issuingDate, "%d-%d", &issuingYear, &issuingMonth); err != nil {
				return fmt.Errorf("invalid issuing date format for certification %s: %w", file.Filename, err)
			}

			// Combine year and month into a single string
			combinedIssuingDate := fmt.Sprintf("%04d-%02d", issuingYear, issuingMonth)

			// Create the certification details entry
			certificationsAndAchievementsDetails := models.StudentCertificationDetailsTable{
				CertificationName: certificationName,   // Store the certification name
				IssuingAuthority:  issuingAuthority,    // Store the issuing authority
				IssuingDate:       combinedIssuingDate, // Store the combined issuing date
				DocumentID:        uint32(certificationDocumentID),
			}

			// Insert the certification details into the database
			if err := tx.Create(&certificationsAndAchievementsDetails).Error; err != nil {
				return fmt.Errorf("failed to create certification entry for %s: %w", file.Filename, err)
			}

			certificationLookupDetails := models.StudentCertificationLookup{
				EnrollmentNo:                  userInput.EnrollmentNo,
				StudentCertificationDetailsID: certificationsAndAchievementsDetails.ID,
			}

			// Insert the certification details into the database
			if err := tx.Create(&certificationLookupDetails).Error; err != nil {
				return fmt.Errorf("failed to create certification lookup entry for %s: %w", file.Filename, err)
			}

		}

		// Create masterEntry enrollment record
		masterEntry := models.EnrollmentMasterLookupTable{
			EnrollmentNo:         userInput.EnrollmentNo,
			LogInDetailsID:       loginDetails.ID,
			AcademicDetailsID:    academicDetails.ID,
			FamilyDetailsID:      familyDetails.ID,
			ProfileDetailsID:     profileDetails.ID,
			ScholarshipDetailsID: scholarshipDetails.ID,
		}

		if err := tx.Create(&masterEntry).Error; err != nil {
			return fmt.Errorf("failed to create master entry of student profile: %w", err)
		}

		return nil
	})

	// Handle any errors during the transaction
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Transaction failed",
			"details": err.Error(),
		})
		return
	}

	// Generate JWT token for student user
	commonTokenString, err := utils.GenerateToken("student_user", "common", 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{
		"message":       "SignUp successfully",
		"Authorization": "Bearer " + commonTokenString,
	})
}

// Assumed file validations are done in the frontend.
// Helper function to upload a file to AWS S3
// func uploadFileToCloud(file *multipart.FileHeader, filename string) (string, error) {
// 	// Open the file
// 	uploadFile, err := file.Open()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open file: %w", err)
// 	}
// 	defer uploadFile.Close()
// 	// Get the S3 service client
// 	client := config.AWSClient
// 	// Upload the file to S3
// 	bucket := config.GetEnv("S3_BUCKET_NAME", "default-bucket-name") // Replace with your default S3 bucket name
// 	key := "uploads/" + filename
// 	uploader := manager.NewUploader(client)
// 	result, err := uploader.Upload(
// 		context.TODO(),
// 		&s3.PutObjectInput{
// 			Bucket: aws.String("my-bucket"),
// 			Key:    aws.String("my-object-key"),
// 			Body:   uploadFile,
// 		})
// 	output, err := u.upload(input)
// 	if err != nil {
// 		var mu manager.MultiUploadFailure
// 		if errors.As(err, &mu) {
// 			// Process error and its associated uploadID
// 			fmt.Println("Error:", mu)
// 			_ = mu.UploadID() // retrieve the associated UploadID
// 		} else {
// 			// Process error generically
// 			fmt.Println("Error:", err.Error())
// 		}
// 		return
// 	}
// 	println(output)
// 	result, err := uploader.Upload(
// 		context.TODO(),
// 		&s3.PutObjectInput{
// 			Bucket: aws.String(bucket),
// 			Key:    aws.String(key),
// 			Body:   uploadFile,
// 		})
// 	if err != nil {
// 		var mu manager.MultiUploadFailure
// 		if errors.As(err, &mu) {
// 			// Process error and its associated uploadID
// 			fmt.Println("Error:", mu)
// 			_ = mu.UploadID() // retrieve the associated UploadID
// 		} else {
// 			// Process error generically
// 			fmt.Println("Error:", err.Error())
// 		}
// 		return "", fmt.Errorf("failed to upload file to S3: %w", err)
// 	}
// 	if err != nil {
// 		return "", fmt.Errorf("failed to upload file to S3: %w", err)
// 	}
// 	// Return the URL of the uploaded file
// 	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
// 	return url, nil
// }

// Helper function to upload a file to AWS S3
func uploadFileToCloud(file *multipart.FileHeader, filename string) (string, error) {

	// Open the file
	uploadFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	} else {
		defer uploadFile.Close()

		// Get the S3 service client
		client := config.AWSClient

		// Retrieve the bucket name from environment variables or use a default
		bucket := os.Getenv("S3_BUCKET_NAME")
		key := "uploads/" + filename

		// Initialize the S3 uploader
		uploader := manager.NewUploader(client)

		// Perform the upload to S3
		result, err := uploader.Upload(
			context.TODO(),
			&s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   uploadFile,
			},
		)
		if err != nil {
			// Handle multi-part upload-specific errors if applicable
			var multiUploadError manager.MultiUploadFailure
			if errors.As(err, &multiUploadError) {
				log.Printf("MultiUpload failure detected, UploadID: %s, Error: %s", multiUploadError.UploadID(), err)
			}
			return "", fmt.Errorf("failed to upload file to S3: %w", err)
		}

		return result.Location, nil
	}
}
