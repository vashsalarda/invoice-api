package route

import (
	"invoice-api/internal/features/revenue/controller"

	"github.com/gofiber/fiber/v2"
)

type InvoiceRoute struct{}

func (c *InvoiceRoute) Init(router *fiber.App) {
	controller := new(controller.RevenueController)
	revenues := router.Group("/revenues")

	revenues.Post("/", controller.CreateRevenue)
	revenues.Get("/", controller.GetAllRevenues)
	revenues.Get("/:id", controller.GetRevenueByID)
	revenues.Patch("/:id", controller.UpdateRevenue)
	revenues.Delete("/:id", controller.DeleteRevenue)
}

	