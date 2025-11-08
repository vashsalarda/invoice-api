package controller

import (
	"invoice-api/internal/features/customer/command"
	"invoice-api/internal/features/customer/model"
	"invoice-api/internal/features/customer/query"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerController struct {
	Command           command.Command
	Query             query.Query
}

func (s *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	s.Command = &command.DefaultCommand{}
	
	payload := new(model.CreateCustomer)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}

	resp, err := s.Command.CreateCustomer(payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}

	return c.Status(201).JSON(resp)
}

func (s *CustomerController) GetAllCustomers(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}

	items, err := s.Query.GetItemsByQuery()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch customers",
		})
	}

	return c.JSON(items)
}

func (s *CustomerController) GetCustomerByID(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	id := c.Params("id")

	item, err := s.Query.GetItemByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch customer",
		})
	}

	return c.JSON(item)
}

func (s *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	s.Command = &command.DefaultCommand{}
	id := c.Params("id")

	payload := new(model.UpdateCustomer)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}
	res, err := s.Command.UpdateCustomer(id, payload)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}

	if res.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Customer updated successfully",
	})
}

func (s *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := s.Command.DeleteCustomer(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}

	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Customer deleted successfully",
	})
}