package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func ApplyOutpass(c fiber.Ctx) error {
	return c.SendString("Do APply logic here")
}
