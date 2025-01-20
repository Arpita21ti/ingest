package middlewares

import (
	"fmt"
	"net/http"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper function to validate the token and extract claims.
func validateToken(c *gin.Context) (map[string]interface{}, error) {
	// Extract token from Authorization header.
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing Authorization header")
	}

	// Validate Bearer format.
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header format")
	}

	tokenString := parts[1]

	// Validate token and extract claims
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

// Token extraction and validation middleware
func TokenValidationMiddleware(c *gin.Context) {
	claims, err := validateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Store the claims in the context for later use in handlers
	c.Set("claims", claims)

	// Proceed to the next handler
	c.Next()
}
