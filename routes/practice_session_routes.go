package routes

import (
	controllersNew "server/controllers/psql"
	// "server/middlewares"
	reqMiddleware "server/middlewares/requests"

	"github.com/gin-gonic/gin"
)

func PracticeSessionRoutes(router *gin.Engine) {
	session := router.Group("/practice-session")
	session.Use(reqMiddleware.RequireJSON()) // Apply the middleware to check and accept for JSON requests only (adds security and reduces load).
	// session.Use(middlewares.TokenValidationMiddleware) // Token validation middleware

	// Practice session routes
	{
		// NOT NEEDED AS IT IS TAKEN CARE OF BY THE QUESTION FETCH HANDLER.
		// session.POST("/start", controllersNew.StartPracticeSessionHandler)
		session.POST("/submit", controllersNew.SubmitPracticeSessionHandler)
		session.POST("/end-forcefully", controllersNew.ForcefullyEndPracticeSessionHandler)
	}
}

// TODO: Update as per this route and models
// Example Requests:

// POST /add_question?type=MCQ
// Content-Type: application/json
// {
//   "question": "What is 2+2?",
//   "options": ["2", "3", "4", "5"],
//   "correct_option": 2,
//   "category": "Math",
//   "difficulty": "Easy",
//   "tags": ["basic", "arithmetic"]
// }

// Fetch all questions of any type:
// GET /questions/aptitude?subcategories=linear&difficulty=Medium&tags=math

// Fetch only MCQs:
// GET /questions/aptitude?subcategories=linear&difficulty=Medium&tags=math&type=MCQ

// Fetch only True/False questions:
// GET /questions/aptitude?subcategories=linear&difficulty=Medium&tags=math&type=TF

// Fetch only Fill-in-the-Blank questions:
// GET /questions/aptitude?subcategories=linear&difficulty=Medium&tags=math&type=FB

// Example Queries

// Fetch all aptitude questions:
// GET /questions/aptitude

// Fetch all quant questions with "probability" tag:
// GET /questions/quant?tags=probability

// Fetch all PYQs of "hard" difficulty:
// GET /questions/pyqs?difficulty=hard

// Fetch typical query
// /questions/quant?tags=algebra&tags=geometry&difficulty=hard&subcategories=linear&subcategories=probability
