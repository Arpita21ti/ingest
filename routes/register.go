package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	AuthRoutes(router)
	PasswordResetRoutes(router)
	QuestionRoutes(router)
	PracticeSessionRoutes(router)
	// Add other route group registrations here...
}

// // Apply middleware to protect routes
// router.GET("/protected", TokenValidationMiddleware, func(c *gin.Context) {
// 	// If the token is valid, the claims will be available
// 	claims, _ := c.Get("claims")
// 	// Process the request with the claims
// 	c.JSON(http.StatusOK, gin.H{"message": "Access granted", "claims": claims})
// })
