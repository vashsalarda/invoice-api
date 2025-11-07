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
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(resp)
}

// HandleGetAllCustomers handles the HTTP request to get all users
func (s *CustomerController) GetAllCustomers(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	users, err := s.Query.GetCustomerByQuery()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// HandleGetCustomerByID handles the HTTP request to get a user by ID
func (s *CustomerController) GetCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.Query.GetCustomerByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

// HandleUpdateCustomer handles the HTTP request to update a user
func (s *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
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
			"error": "Failed to update user",
		})
	}

	if res.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Customer updated successfully",
	})
}

// HandleDeleteCustomer handles the HTTP request to delete a user
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
			"error": "Failed to update user",
		})
	}

	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Customer deleted successfully",
	})
}