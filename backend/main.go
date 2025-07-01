package main

import (
	"log"

	"backend/config"
	"backend/handlers"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func main() {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env:", err)
	}

	// Setup router
	r := gin.Default()

	// Routes
	r.GET("/trade-ideas", handlers.GetTradeIdeas)

	// Start server
	log.Println("🚀 Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server error:", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all origins for development
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Define your routes here
	router.GET("/trade-ideas", handlers.GetTradeIdeas)

	// Listen on all interfaces so your phone or emulator can reach it
	router.Run("0.0.0.0:8080")
}
