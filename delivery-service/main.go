package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"delivery-service/handlers"
	"delivery-service/middleware"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to PostgreSQL database
	dbURL := os.Getenv("DELIVERY_DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Initialize Echo
	e := echo.New()
	e.Use(emiddleware.CORS())

	// Group with JWT middleware
	deliveryGroup := e.Group("/deliveries")
	deliveryGroup.Use(echo.WrapMiddleware(middleware.JWTAuthMiddleware))

	// Route handlers
	deliveryGroup.GET("/assigned", handlers.GetAssignedDeliveries(db))
	deliveryGroup.PUT("/status", handlers.UpdateDeliveryStatus(db))
	deliveryGroup.PUT("/return/:delivery_id", handlers.RaiseReturnRequest(db))

	// Start server
	port := os.Getenv("PORT")

	e.Logger.Fatal(e.Start(":" + port))
}
