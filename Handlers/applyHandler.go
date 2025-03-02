package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/senthan-07/outpassBE/models"
	"gorm.io/gorm"
)

// OutpassRequest defines the expected request structure.
type OutpassRequest struct {
	StudentID    uint64 `json:"student_id" form:"student_id"`
	OutpassType  string `json:"outpass_type" form:"outpass_type"`
	Status       string `json:"status" form:"status"`
	ValidFrom    string `json:"valid_from" form:"valid_from"`
	ValidUntil   string `json:"valid_until" form:"valid_until"`
	ApprovedByID uint64 `json:"approved_by_id" form:"approved_by_id"`
}

// ApplyOutpass handles student outpass applications.
func ApplyOutpass(c fiber.Ctx) error {
	// Get the database instance from Fiber context.
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}

	// Parse the request body into OutpassRequest.
	var request OutpassRequest
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Validate Student ID.
	var student models.Student
	if err := db.First(&student, request.StudentID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Student ID"})
	}

	// Validate Approver exists (Teacher or Warden).
	var approverExists bool
	var teacher models.Teacher
	var warden models.Warden
	if db.First(&teacher, request.ApprovedByID).Error == nil {
		approverExists = true
	} else if db.First(&warden, request.ApprovedByID).Error == nil {
		approverExists = true
	}
	if !approverExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Approver ID"})
	}

	// Convert valid_from and valid_until to time.Time.
	validFrom, err := time.Parse(time.RFC3339, request.ValidFrom)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_from format"})
	}
	validUntil, err := time.Parse(time.RFC3339, request.ValidUntil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_until format"})
	}

	// Create and insert the Outpass record.
	outpass := models.Outpass{
		StudentID:    request.StudentID,
		OutpassType:  request.OutpassType,
		Status:       request.Status,
		ValidFrom:    validFrom,
		ValidUntil:   validUntil,
		ApprovedByID: request.ApprovedByID,
	}

	if err := db.Create(&outpass).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Outpass"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Outpass created successfully", "outpass": outpass})
}
