package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	models "github.com/senthan-07/outpassBE/Models"
	"gorm.io/gorm"
)

// ApproveOutpass handles approval of an outpass by a Warden or Teacher.
func ApproveOutpass(c fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}

	outpassIDStr := c.Params("outpass_id")
	fmt.Println("Raw outpass_id:", outpassIDStr)
	// Get the outpass ID from URL params
	outpassID, err := strconv.ParseUint(c.Params("outpass_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outpass ID"})
	}

	// Debug: Print outpassID
	fmt.Printf("Parsed Outpass ID: %d\n", outpassID)

	// Fetch outpass from database
	var outpass models.Outpass
	if err := db.First(&outpass, uint(outpassID)).Error; err != nil {
		fmt.Println("Outpass not found in DB")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Outpass not found"})
	}

	fmt.Println("Outpass found:", outpass.ID)

	// Parse the request body
	var req struct {
		ApproverID   uint64 `json:"approver_id"`
		ApproverType string `json:"approver_type"`
		Status       string `json:"status"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Check if the approver exists
	if req.ApproverType == "Warden" {
		var warden models.Warden
		if err := db.First(&warden, req.ApproverID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":         "Invalid Approver ID",
				"approver_id":   req.ApproverID,
				"approver_type": req.ApproverType,
			})
		}
	} else if req.ApproverType == "Teacher" {
		var teacher models.Teacher
		if err := db.First(&teacher, req.ApproverID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":         "Invalid Approver ID",
				"approver_id":   req.ApproverID,
				"approver_type": req.ApproverType,
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Approver Type"})
	}

	// Update the outpass status
	outpass.Status = req.Status
	outpass.ApprovedByID = &req.ApproverID
	outpass.ApproverType = req.ApproverType

	if err := db.Save(&outpass).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update outpass status"})
	}

	// Send notification to the student
	notification := models.Notification{
		UserID:  outpass.StudentID,
		Message: "Your outpass request has been " + req.Status + " by " + req.ApproverType,
		Read:    false,
	}
	db.Create(&notification)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Outpass approved successfully",
		"outpass": outpass,
	})
}

func SetRegularOutpassDates(c fiber.Ctx) error {
	return c.SendString("Do set date logic ")
}

func ApproveRejectOutpass(c fiber.Ctx) error {
	return c.SendString("Do set date logic ")
}

func ValidateOutpass(c fiber.Ctx) error {
	return c.SendString("Do set date logic ")
}
