package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/heshanu/go/handlers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Gin
	r := gin.Default()
	
	// Serve static files
	r.Static("/static", "./static")
	
	// Load templates
	r.LoadHTMLGlob("templates/*")
	
	// Routes
	r.GET("/", homeHandler)
	r.POST("/generate", handlers.GenerateTextHandler)
	r.POST("/classify", handlers.ClassifyImageHandler)
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server running on :%s", port)
	r.Run(":" + port)
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}