package controller

import (
	"invoice-api/internal/features/invoice/command"
	"invoice-api/internal/features/invoice/model"
	"invoice-api/internal/features/invoice/query"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceController struct {
	Command command.InvoiceCommand
	Query   query.InvoiceQuery
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
			"error": "Failed to create invoice. " + err.Error(),
		})
	}

	return c.Status(201).JSON(resp)
}

func (s *InvoiceController) GetAllInvoices(c *fiber.Ctx) error {
	s.Query = &query.DefaultInvoiceQuery{}
	keyword := c.Query("keyword")
	status := c.Query("status")
	sizeStr := c.Query("size")
	pageStr := c.Query("page")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 25
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	items, err := s.Query.GetItemsByQuery(keyword, status, size, page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(items)
}

func (s *InvoiceController) GetInvoiceByID(c *fiber.Ctx) error {
	s.Query = &query.DefaultInvoiceQuery{}
	id := c.Params("id")

	item, err := s.Query.GetItemByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch invoice",
		})
	}

	return c.JSON(item)
}

func (s *InvoiceController) UpdateInvoice(c *fiber.Ctx) error {
	s.Command = &command.DefaultInvoiceCommand{}
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
			"error": "Failed to update invoice",
		})
	}

	if res.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to update invoice",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice updated successfully",
	})
}

func (s *InvoiceController) DeleteInvoice(c *fiber.Ctx) error {
	s.Command = &command.DefaultInvoiceCommand{}
	id := c.Params("id")

	res, err := s.Command.DeleteItem(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update invoice",
		})
	}

	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to delete invoice",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}

func (s *InvoiceController) GetLatestInvoices(c *fiber.Ctx) error {
	s.Query = &query.DefaultInvoiceQuery{}
	res, err := s.Query.GetLatestInvoices()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch latest invoices",
		})
	}

	return c.JSON(res)
}

func (s *InvoiceController) GetTotalInvoices(c *fiber.Ctx) error {
	s.Query = &query.DefaultInvoiceQuery{}
	keyword := c.Query("keyword")
	status := c.Query("status")
	validStatus := map[string]bool{"pending": true, "paid": true, "cancelled": true, "void": true}
	if status != "" {
		status_check := strings.ToLower(status)
		if !validStatus[status_check] {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid status. Please select from `pending`, `paid`, `cancelled`, `void`",
			})
		}
	}
	items, err := s.Query.GetTotalItemsByQuery(keyword, status)
	if err != nil {
		return c.JSON(0)
	}

	return c.JSON(items)
}

// func (s *InvoiceController) GetCustomerInvoices(c *fiber.Ctx) error {
// 	s.Query = &query.DefaultInvoiceQuery{}
// 	keyword := c.Query("keyword")
// 	items, err := s.Query.GetCustomersInvoices(keyword)
// 	if err != nil {
// 		return c.JSON(0)
// 	}

// 	return c.JSON(items)
// }
