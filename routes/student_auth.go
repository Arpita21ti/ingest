package routes

import (
	"server/controllers"
	controllersNew "server/controllers/psql"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.StudentSignupHandler)
		auth.POST("/login", controllers.StudentLoginHandler)
		auth.POST("/signup-new", controllersNew.StudentSignupHandler)
		auth.POST("/login-new", controllersNew.StudentLoginHandler)
	}
}

// Example Requests

// For Login
// {
// 	"enrollmentNo": "123",
// 	"email": "hsdajh@a.com",
// 	"password": "123456789",
// 	"phone": "1234567890"
//   }

// For Sign Up
// {
// 	"enrollmentNo": "123",
// 	"email": "abc@ag.com",
// 	"password": "1",
// 	"phone": "1234567890",
// 	"name": "aradhya",
// 	"branch": "aiml",
// 	"year_of_enrollment": 2022,
// 	"10th_percentage":100,
// 	"12th_percentage":100,
// 	"previous_sem_cgpa":9.8,
// 	"previous_sem_sgpa":9.9,
// 	"certifications":["https://abc.abc"],
// 	"resume":"https://abc.abc"
//   }
