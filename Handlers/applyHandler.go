package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	models "github.com/senthan-07/outpassBE/Models"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// sendEmail sends an email notification
func sendEmail(to string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_USER")) // Use configured SMTP email
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	// SMTP Dialer Configuration
	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587, // Default SMTP port for TLS
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	return dialer.DialAndSend(mailer)
}

// OutpassRequest defines the expected request structure.
type OutpassRequest struct {
	StudentID    uint64 `json:"student_id"`
	OutpassType  string `json:"outpass_type"`
	ValidFrom    string `json:"valid_from"`
	ValidUntil   string `json:"valid_until"`
	ApproverType string `json:"approver_type"` // "Warden" or "Teacher"
	ApproverName string `json:"approver_name"` // Approver's name
}

func parseDate(dateStr string) (time.Time, error) {
	// Try parsing in RFC3339 format first
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return parsedTime, nil
	}

	// If RFC3339 fails, try parsing in "YYYY-MM-DD" format
	return time.Parse("2006-01-02", dateStr)
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
	validFrom, err := parseDate(request.ValidFrom)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_from format"})
	}
	validUntil, err := parseDate(request.ValidUntil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid valid_until format"})
	}

	// Find the Approver in the correct table (Warden or Teacher)
	var approverEmail string
	var approverID uint64
	var approverName string

	if request.ApproverType == "Warden" {
		var warden models.Warden
		if err := db.Where("name = ?", request.ApproverName).First(&warden).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Warden not found"})
		}
		approverEmail = warden.Email
		approverID = uint64(warden.ID)
		approverName = warden.Name
	} else if request.ApproverType == "Teacher" {
		var teacher models.Teacher
		if err := db.Where("name = ?", request.ApproverName).First(&teacher).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Teacher not found"})
		}
		approverEmail = teacher.Email
		approverID = uint64(teacher.ID)
		approverName = teacher.Name
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid approver type"})
	}

	// Create the Outpass record (Approver is nil since it's pending).
	outpass := models.Outpass{
		StudentID:    request.StudentID,
		OutpassType:  request.OutpassType,
		Status:       "Pending",
		ValidFrom:    validFrom,
		ValidUntil:   validUntil,
		ApprovedByID: nil,
		ApproverType: request.ApproverType,
		ApproverName: request.ApproverName,
	}

	if err := db.Create(&outpass).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Outpass"})
	}

	// **Send Email to Approver**
	emailSubject := "New Outpass Request"
	emailBody := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>A new outpass request has been submitted by <b>%s</b>.</p>
		<p><b>Details:</b></p>
		<ul>
			<li><b>Outpass Type:</b> %s</li>
			<li><b>Valid From:</b> %s</li>
			<li><b>Valid Until:</b> %s</li>
		</ul>
		<p>Please review the request.</p>
	`, approverName, student.Name, request.OutpassType, request.ValidFrom, request.ValidUntil)

	go func() {
		if err := sendEmail(approverEmail, emailSubject, emailBody); err != nil {
			fmt.Println("Failed to send email:", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Outpass request submitted successfully",
		"outpass": outpass,
		"notified": fiber.Map{
			"user_id": approverID,
			"name":    approverName,
			"type":    request.ApproverType,
		},
	})
}
