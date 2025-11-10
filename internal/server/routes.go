package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.InvoiceAPIHandler)
	s.App.Get("/health", s.HealthHandler)
}

func (s *FiberServer) InvoiceAPIHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Invoice API",
	}

	return c.JSON(resp)
}

func (s *FiberServer) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
