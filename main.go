// To access any route other than signUp or login send Authorization: Bearer <token>
// in the packet header
package main

import (
	"log"
	"os"

	"server/config"
	"server/routes"
	seed "server/seeds"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Load .env variables
	if err := config.LoadEnvVariables(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// // Initialize AWS session
	// if err := config.InitializeAWSSession(); err != nil {
	// 	log.Fatalf("Error initializing AWS session: %v", err)
	// }

	// Connect to MongoDB
	if err := config.ConnectMongoDB(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Connect to PostgresDB
	if err := config.ConnectPostgresDB(); err != nil {
		log.Fatalf("Error connecting to the PostgresDB: %v", err)
	}

	// Check if the environment is development, then seed the database
	// Did seeding here inplace of config file to avoid cyclic imports in the config and seed files.
	env := os.Getenv("USE_SEED_DATA")
	if env == "true" {
		if err := seed.SeedQuestions(); err != nil {
			log.Fatalf("Error seeding the database: %v", err)
		}
		log.Println("Test questions seeded successfully.")
	} else {
		log.Println("No Seeding of data done.")
	}

	go InitGraphQLServer()

	// Initialize Gin router
	router := gin.Default()

	// // Initialize Gin router without default middlewares
	// router := gin.New()

	// // Manually attach Logger and Recovery middleware
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// // Set up GraphQL server
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	// router.POST("/query", gin.WrapH(srv))
	// router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},           // Replace * with specific origins for production
		AllowMethods:     []string{"GET", "POST"}, // "PUT", "DELETE"
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Default route to show server running message
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Go TNP RGPV Server is running",
		})
	})

	// Register all routes
	routes.RegisterRoutes(router)

	// Get the PORT from the environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if no PORT environment variable is set
	}

	// Start the server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Graceful shutdown: Disconnect MongoDB and PostgreSQL when the server stops
	defer func() {
		if err := config.DisconnectMongoDB(); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		// Disconnect from PostgreSQL
		if err := config.DisconnectPostgresDB(); err != nil {
			log.Fatalf("Error disconnecting from PostgreSQL: %v", err)
		}
	}()
}
