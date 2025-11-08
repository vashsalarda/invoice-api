package route

import (
	"invoice-api/internal/features/invoice/controller"

	"github.com/gofiber/fiber/v2"
)

type InvoiceRoute struct{}

func (c *InvoiceRoute) Init(router *fiber.App) {
	controller := new(controller.InvoiceController)
	invoices := router.Group("/invoices")

	invoices.Post("/", controller.CreateInvoice)
	invoices.Get("/", controller.GetAllInvoices)
	invoices.Get("/latest", controller.GetLatestInvoices)
	invoices.Get("/:id", controller.GetInvoiceByID)
	invoices.Patch("/:id", controller.UpdateInvoice)
	invoices.Delete("/:id", controller.DeleteInvoice)
}

	