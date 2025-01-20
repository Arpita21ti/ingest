package middlewares

import (
	"fmt"
	"net/http"
	"server/config"
	models "server/models/common"

	"github.com/gin-gonic/gin"
)

func CheckAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userID from request header.
		userID := c.GetHeader("UserID")
		// Extract documentID from Query Parameters of the request.
		documentID := c.Query("documentID")

		if userID == "" || documentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userID or documentID"})
			c.Abort()
			return
		}

		// Get the current DB connection
		db := config.GetPostgresDBConnection()

		var privilege models.DocumentPrivilegeTable

		// Check if the document is restricted by any user (documentID exists in the table)
		result := db.First(&privilege, "document_id = ? AND user_id = ?", documentID, userID)

		// If no matching row is found, grant access (document not restricted for this user)
		if result.RowsAffected == 0 {
			fmt.Println("Access granted: Document not restricted or no specific access required.")
			c.Next()
			return
		}

		// If a row is found, check the specific access permissions
		if !privilege.CanRead { // You can customize this condition for other operations (CanWrite, CanDelete)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: User does not have permissions for this document"})
			c.Abort()
			return
		}

		// Access is allowed
		fmt.Println("Access granted: User has the required permissions.")
		c.Next()
	}
}
