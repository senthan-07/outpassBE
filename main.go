package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	models "github.com/senthan-07/outpassBE/Models"
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

	log.Println("Database connected successfully!")

	// Check if tables exist before migrating
	if db.Migrator().HasTable(&models.Student{}) &&
		db.Migrator().HasTable(&models.Teacher{}) &&
		db.Migrator().HasTable(&models.Warden{}) &&
		db.Migrator().HasTable(&models.Outpass{}) &&
		db.Migrator().HasTable(&models.User{}) {
		log.Println("Tables already exist, skipping migration.")
	} else {
		errMigrate := db.AutoMigrate(&models.Student{}, &models.Teacher{}, &models.Warden{}, &models.Outpass{}, &models.User{})
		if errMigrate != nil {
			log.Fatal("Migration failed:", errMigrate)
		} else {
			log.Println("Migration successful!")
		}
	}

	// Insert initial records
	seedDatabase(db)
}

func seedDatabase(db *gorm.DB) {
	// Seed Students
	if err := db.First(&models.Student{}, 1).Error; err != nil {
		db.Create(&models.Student{
			ID:       1,
			Name:     "John Doe",
			Email:    "senthanbalu@gmail.com",
			Password: "hashedpassword",
		})
		log.Println("Inserted test student")
	}

	// Seed Teachers
	if err := db.First(&models.Teacher{}, 1).Error; err != nil {
		db.Create(&models.Teacher{
			ID:       1,
			Name:     "Dr. Smith",
			Email:    "senthanbalu@gmail.com",
			Password: "hashedpassword",
		})
		log.Println("Inserted test teacher")
	}

	// Seed Wardens
	if err := db.First(&models.Warden{}, 1).Error; err != nil {
		db.Create(&models.Warden{
			ID:       1,
			Name:     "Mr. Warden",
			Email:    "senthanbalu@gmail.com",
			Password: "hashedpassword",
		})
		log.Println("Inserted test warden")
	}

	// Seed Users (Approvers)
	if err := db.First(&models.User{}, 1).Error; err != nil {
		db.Create(&models.User{
			ID:    1,
			Name:  "Admin",
			Email: "admin@university.com",
		})
		log.Println("Inserted test user")
	}
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
