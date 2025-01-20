package routes

import (
	reqMiddleware "server/middlewares/requests"

	"github.com/gin-gonic/gin"
)

func LeaderboardRoutes(router *gin.Engine) {
	leaderboard := router.Group("/leaderboard")
	// Apply middleware to check and accept for JSON requests only (adds security and reduces load).
	leaderboard.Use(reqMiddleware.RequireJSON())

	// leaderboard.GET("/", controllersNew.GetLeaderboard)
	// leaderboard.GET("/:id", controllersNew.GetLeaderboardByID)
	// leaderboard.POST("/", controllersNew.AddLeaderboard)
	// leaderboard.PUT("/:id", controllersNew.UpdateLeaderboard)
	// leaderboard.DELETE("/:id", controllersNew.DeleteLeaderboard)
}

// Example Requests:

// POST /leaderboard
// Content-Type: application/json
// {
//   "rank": 1,
//   "score": 100.0,
//   "domain": "Math",
//   "subDomain": "Algebra"
// }

// GET /leaderboard
// Fetch all leaderboard records

// GET /leaderboard/:id
// Fetch a specific leaderboard record by ID

// PUT /leaderboard/:id
// Content-Type: application/json
// {
//   "rank": 1,
//   "score": 100.0,
//   "domain": "Math",
//   "subDomain": "Algebra"
// }
// Update a specific leaderboard record by ID

// DELETE /leaderboard/:id
// Delete a specific leaderboard record by ID
