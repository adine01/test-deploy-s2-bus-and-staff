package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	if err := InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer CloseDB()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Initialize routes
	setupRoutes(router)

	// Get port from environment or default to 8082
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Bus Staff Assignment Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine) {
	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "bus-staff-assignment"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Assignment routes
		api.POST("/assignments", handleCreateAssignment)
		api.GET("/assignments", handleGetAssignments)
		api.GET("/assignments/:id", handleGetAssignment)
		api.PUT("/assignments/:id", handleUpdateAssignment)
		api.DELETE("/assignments/:id", handleDeleteAssignment)

		// Query routes
		api.GET("/assignments/bus/:busId", handleGetStaffForBus)
		api.GET("/assignments/staff/:staffId", handleGetAssignmentsForStaff)
	}
}
