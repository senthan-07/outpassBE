package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	models "github.com/senthan-07/outpassBE/Models"
	"gorm.io/gorm"
)

// ApproveOutpass handles approval or rejection of an outpass by a Warden or Teacher.
func ApproveOutpass(c fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}

	// Get the outpass ID from URL params
	outpassID, err := strconv.ParseUint(c.Params("outpass_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outpass ID"})
	}

	// Fetch the outpass record
	var outpass models.Outpass
	if err := db.First(&outpass, uint(outpassID)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Outpass not found"})
	}

	// Parse the request body
	var req struct {
		ApproverName string `json:"approver_name"` // Approver identified by Name
		ApproverType string `json:"approver_type"` // "Warden" or "Teacher"
		Status       string `json:"status"`        // "Approved" or "Rejected"
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Validate Approver Exists
	if req.ApproverType == "Warden" {
		var warden models.Warden
		if err := db.Where("name = ?", req.ApproverName).First(&warden).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Warden not found"})
		}
	} else if req.ApproverType == "Teacher" {
		var teacher models.Teacher
		if err := db.Where("name = ?", req.ApproverName).First(&teacher).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Teacher not found"})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Approver Type"})
	}

	// Update the outpass status
	outpass.Status = req.Status
	outpass.ApproverName = req.ApproverName
	outpass.ApproverType = req.ApproverType

	if err := db.Save(&outpass).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update outpass status"})
	}

	// Send Email Notification to Student
	var student models.Student
	if err := db.First(&student, outpass.StudentID).Error; err == nil {
		emailSubject := fmt.Sprintf("Your Outpass Request has been %s", req.Status)
		emailBody := fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your outpass request has been <b>%s</b> by <b>%s (%s)</b>.</p>
			<p><b>Details:</b></p>
			<ul>
				<li><b>Outpass Type:</b> %s</li>
				<li><b>Valid From:</b> %s</li>
				<li><b>Valid Until:</b> %s</li>
			</ul>
			<p>Thank you.</p>
		`, student.Name, req.Status, req.ApproverName, req.ApproverType, outpass.OutpassType, outpass.ValidFrom, outpass.ValidUntil)

		// Send email using common system email (not approver's email)
		go sendEmail(student.Email, emailSubject, emailBody)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Outpass status updated successfully",
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
