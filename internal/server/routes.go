package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

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
