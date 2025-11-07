package route

import (
	"invoice-api/internal/features/customer/controller"

	"github.com/gofiber/fiber/v2"
)

type CustomerRoute struct{}

func (c *CustomerRoute) Init(router *fiber.App) {
	controller := new(controller.CustomerController)
	customers := router.Group("/customers")

	customers.Post("/", controller.CreateCustomer)
	customers.Get("/", controller.GetAllCustomers)
	customers.Get("/:id", controller.GetAllCustomers)
	customers.Patch("/:id", controller.UpdateCustomer)
	customers.Delete("/:id", controller.DeleteCustomer)
}
