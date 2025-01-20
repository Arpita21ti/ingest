package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func PasswordResetRoutes(router *gin.Engine) {
	reset := router.Group("/auth") // In "/auth" group as routes here use no middlewares.
	{
		reset.POST("/forgot-password", controllers.ForgotPasswordHandler)
		reset.POST("/reset-password", controllers.ResetPasswordHandler)
	}
}
