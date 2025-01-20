package routes

import (
	"server/controllers"
	controllersNew "server/controllers/psql"

	// "server/middlewares"
	reqMiddleware "server/middlewares/requests"

	"github.com/gin-gonic/gin"
)

func QuestionRoutes(router *gin.Engine) {
	questions := router.Group("/questions")
	questions.Use(reqMiddleware.RequireJSON()) // Apply the middleware to check and accept for JSON requests only (adds security and reduces load).
	// questions.Use(middlewares.TokenValidationMiddleware) // Token validation middleware

	// Question routes
	{
		// Question routes
		questions.GET("/hierarchy", controllersNew.GetQuestionHierarchy)
		questions.GET("/domains", controllersNew.GetDomains)
		questions.GET("/subdomains/:domainID", controllersNew.GetSubDomains)
		questions.GET("/niches/:subDomainsID", controllersNew.GetNiches)
		questions.GET("/difficulty-levels/:nicheID", controllersNew.GetDifficultyLevels)
		questions.GET("/formats/:difficultyLevelID", controllersNew.GetFormats)

		questions.POST("/fetch", controllersNew.GetQuestions)

		// Endpoint to add single/individual question.
		questions.POST(
			"/add-question",
			// middlewares.PrivilegedMiddleware("admin"), // Privileges check for "admin"
			controllers.AddSingleQuestionHandler,
		)

		// Endpoint to add bulk/multiple questions in a go.
		questions.POST(
			"/add-bulk-questions",
			// middlewares.PrivilegedMiddleware("admin"), // Privileges check for "admin"
			controllers.AddBulkQuestionHandler,
		)
	}
}

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
