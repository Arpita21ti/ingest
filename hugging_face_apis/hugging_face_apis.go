package huggingFace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CallHuggingFaceAPI(question string) (string, error) {
	// Get AI Model API from environment variables
	huggingFaceAPIKey := os.Getenv("HUGGING_FACE_API_KEY")

	url := "https://api-inference.huggingface.co/models/google/gemma-2-2b-it"
	payload := map[string]string{"inputs": question}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", huggingFaceAPIKey))
	req.Header.Set("Content-Type", "application/json")

	// Print the request details for debugging or logging
	fmt.Printf("Request URL: %s\n", req.URL)
	fmt.Printf("Request Method: %s\n", req.Method)
	fmt.Println("Request Headers:")
	for name, values := range req.Header {
		for _, value := range values {
			fmt.Printf("\t%s: %s\n", name, value)
		}
	}
	fmt.Println("Request Body:")
	fmt.Println(string(payloadBytes))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Failed to make API request", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"project/models"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// const API_URL = "https://api-inference.huggingface.co/models/google/gemma-2-2b-it"
// const API_KEY = "Bearer hf_***" // Replace with your Hugging Face API key

// type HuggingFacePayload struct {
// 	Inputs string `json:"inputs"`
// }

// type HuggingFaceResponse struct {
// 	GeneratedText string `json:"generated_text"`
// }

// // GenerateQuestionHandler handles the AI processing and database saving
// func GenerateQuestionHandler(c *gin.Context, db *gorm.DB) {
// 	var input struct {
// 		Text string `json:"text" binding:"required"`
// 	}

// 	// Parse JSON input
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Prepare API payload
// 	payload := HuggingFacePayload{Inputs: input.Text}
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process input"})
// 		return
// 	}

// 	// Make the API call
// 	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API request"})
// 		return
// 	}
// 	req.Header.Set("Authorization", API_KEY)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make API request"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Parse API response
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
// 		return
// 	}

// 	var hfResponse []HuggingFaceResponse
// 	if err := json.Unmarshal(body, &hfResponse); err != nil || len(hfResponse) == 0 {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
// 		return
// 	}

// 	// Format and save question
// 	formattedQuestion := models.Question{
// 		Type:        "MCQ",  // Example: set default type; adjust as needed
// 		Category:    "General",
// 		Subcategory: "AI",
// 		Difficulty:  "Medium",
// 		Tags:        "example,ai",
// 		Status:      "Draft",
// 		Text:        hfResponse[0].GeneratedText,
// 		Option:      "",
// 		Answer:      "",
// 		Explanation: "",
// 	}

// 	if err := db.Create(&formattedQuestion).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question"})
// 		return
// 	}

// 	// Respond with the saved question
// 	c.JSON(http.StatusOK, formattedQuestion)
// }
