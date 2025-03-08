package router

import (
	"github.com/gofiber/fiber/v3"
	handlers "github.com/senthan-07/outpassBE/Handlers"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, ddb *gorm.DB) {
	app.Get("/auth", handlers.GetAuth)
	app.Post("/student/outpass/apply", handlers.ApplyOutpass)
	app.Post("/warden/outpass/regular/dates", handlers.SetRegularOutpassDates) //set dates for warden
	app.Put("/warden/outpass/approve/:outpass_id", handlers.ApproveOutpass)    // PUT to approve/reject outpass by ID
	app.Put("/teacher/outpass/response/:id", handlers.ApproveRejectOutpass)    // outpass approval/rejection (with student outpass status)
	app.Post("/outpass/validate/:id", handlers.ValidateOutpass)                // POST to validate outpass by ID
}
