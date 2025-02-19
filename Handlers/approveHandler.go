package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func ApproveOutpass(c fiber.Ctx) error {
	return c.SendString("Do approve logic ")
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
