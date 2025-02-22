package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	models "github.com/senthan-07/outpassBE/Models"
	router "github.com/senthan-07/outpassBE/Routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Successfully connected to the database and migrated models!")
	checkData()
}

func checkData() {
	layout := "2006-01-02 15:04:05"
	validFrom, err := time.Parse(layout, "2025-02-22 09:00:00")
	if err != nil {
		log.Fatal("Error parsing validFrom:", err)
	}

	validUntil, err := time.Parse(layout, "2025-02-22 18:00:00")
	if err != nil {
		log.Fatal("Error parsing validUntil:", err)
	}
	warden := models.Warden{
		Name:     "Warden John",
		Email:    "warden.john@example.com",
		Password: "password123",
	}
	db.Create(&warden)

	teacher := models.Teacher{
		Name:     "Teacher Sarah",
		Email:    "teacher.sarah@example.com",
		Password: "password123",
	}
	db.Create(&teacher)

	student := models.Student{
		Name:     "Student Alex",
		Email:    "alex.student@example.com",
		Password: "student123",
	}
	db.Create(&student)

	outpass := models.Outpass{
		StudentID:   student.ID,
		OutpassType: "Regular",
		Status:      "Pending",
		ValidFrom:   validFrom,
		ValidUntil:  validUntil,
	}
	db.Create(&outpass)

	log.Println("Sample data inserted successfully!")
}

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	router.SetupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
