package routes

import (
	"users-api/internal/handler"
	"users-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, h *handler.UserHandler) {
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	u := app.Group("/users")
	u.Post("/", h.CreateUser)
	u.Get("/", h.ListUsers)
	u.Get("/:id", h.GetUser)
	u.Put("/:id", h.UpdateUser)
	u.Delete("/:id", h.DeleteUser)
}
