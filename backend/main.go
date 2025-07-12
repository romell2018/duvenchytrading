package main

import (
	"log"
	"os"

	"backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api")
	{
		api.GET("/quote/:symbol", handlers.GetQuote)
		api.GET("/signal/:symbol", handlers.GetSignal)
		api.GET("/news/:symbol", handlers.GetNewsBias)
		api.GET("/rrr/:symbol", handlers.GetRRR)
		api.GET("/topdown", handlers.GetTopDown)
		r.GET("/ai/signal/:symbol", handlers.GetAISignal)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
