package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	router "github.com/senthan-07/outpassBE/Routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB // Global variable

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	var errDB error
	db, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Use global db variable
	if errDB != nil {
		log.Fatal("Failed to connect to the database:", errDB)
	}

	log.Println("Database connected successfully!") // Just confirmation

}

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Use(func(c fiber.Ctx) error {
		c.Locals("db", db) // Store global db instance in request context
		return c.Next()
	})
	router.SetupRoutes(app, db) // Pass the global db instance

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
