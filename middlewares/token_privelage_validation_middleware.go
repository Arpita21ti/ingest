package middlewares

import (
	"net/http"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// PrivilegedMiddleware checks the user's privileges.
func PrivilegedMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := validateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Check if the user has the required privileges
		if !utils.HasPrivilege(claims, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
			c.Abort()
			return
		}

		// User has sufficient privileges, proceed to the next handler
		c.Next()
	}
}
