package controller

import (
	"invoice-api/internal/features/invoice/command"
	"invoice-api/internal/features/invoice/model"
	"invoice-api/internal/features/invoice/query"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceController struct {
	Command           command.InvoiceCommand
	Query             query.InvoiceQuery
}

func (s *InvoiceController) CreateInvoice(c *fiber.Ctx) error {
	s.Command = &command.DefaultInvoiceCommand{}
	
	payload := new(model.CreateInvoice)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}

	resp, err := s.Command.CreateItem(payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(resp)
}

func (s *InvoiceController) GetAllInvoices(c *fiber.Ctx) error {
	s.Query = &query.DefaultInvoiceQuery{}
	users, err := s.Query.GetItemsByQuery()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

func (s *InvoiceController) GetInvoiceByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.Query.GetItemByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

func (s *InvoiceController) UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	payload := new(model.UpdateInvoice)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}
	res, err := s.Command.UpdateItem(id, payload)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	if res.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice updated successfully",
	})
}

func (s *InvoiceController) DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := s.Command.DeleteItem(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}

func (s *InvoiceController) GetLatestInvoices(c *fiber.Ctx) error {
	res, err := s.Query.GetLatestInvoices()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch latest invoices",
		})
	}

	return c.JSON(res)
}
