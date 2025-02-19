package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func GetAuth(c fiber.Ctx) error {
	return c.SendString("Do Auth Logic ")
}
