// Tehe middlware to check and allow only JSON requests on specified routes.
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")

		if contentType == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				// TODO: Replace the error message (after development and testing is completed) with:
				// "error": "" Something that is not telling the mistake and not easily understandable
				"error": "Content-Type header is required.",
			})
			c.Abort()
			return
		}

		// Check the Content-Type header
		if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json;") {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				// TODO: Replace the error message (after development and testing is completed) with:
				// "error": "Incorrect request type."
				"error": "This endpoint only accepts JSON requests.",
			})
			c.Abort() // Stop processing further middleware or the handler
			return
		}
		c.Next() // Continue to the next middleware or handler
	}
}
