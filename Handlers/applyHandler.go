package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	models "github.com/senthan-07/outpassBE/Models"
	"gorm.io/gorm"
)

// OutpassRequest defines the expected request structure.
type OutpassRequest struct {
	StudentID   uint64 `json:"student_id"`
	OutpassType string `json:"outpass_type"`
	ValidFrom   string `json:"valid_from"`
	ValidUntil  string `json:"valid_until"`
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

	// Convert valid_from and valid_until to time.Time.
	validFrom, err := time.Parse(time.RFC3339, request.ValidFrom)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_from format"})
	}
	validUntil, err := time.Parse(time.RFC3339, request.ValidUntil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_until format"})
	}

	// Create the Outpass record (Approver is nil since it's pending).
	outpass := models.Outpass{
		StudentID:    request.StudentID,
		OutpassType:  request.OutpassType,
		Status:       "Pending",
		ValidFrom:    validFrom,
		ValidUntil:   validUntil,
		ApprovedByID: nil, // No approver yet
		ApproverType: "",  // Empty until approval
	}

	if err := db.Create(&outpass).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Outpass"})
	}

	// **Determine Approver (Warden for Regular, Teacher for Special/Emergency)**
	var approverID uint64
	var approverName string
	var approverType string

	if request.OutpassType == "Regular" {
		// Assign a Warden (Dummy Logic: Get first warden)
		var warden models.Warden
		if err := db.First(&warden).Error; err == nil {
			approverID = uint64(warden.ID)
			approverName = warden.Name
			approverType = "Warden"
		}
	} else {
		// Assign a Teacher (Dummy Logic: Get first teacher)
		var teacher models.Teacher
		if err := db.First(&teacher).Error; err == nil {
			approverID = uint64(teacher.ID)
			approverName = teacher.Name
			approverType = "Teacher"
		}
	}

	// **Send Notification to Approver**
	if approverID != 0 {
		notification := models.Notification{
			UserID:  approverID,
			Message: "New outpass request from " + student.Name,
			Read:    false,
		}
		db.Create(&notification)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Outpass request submitted successfully",
		"outpass": outpass,
		"notified": fiber.Map{
			"user_id": approverID,
			"name":    approverName,
			"type":    approverType,
		},
	})
}
