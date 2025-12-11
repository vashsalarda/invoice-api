package controller

import (
	"invoice-api/internal/features/customer/command"
	"invoice-api/internal/features/customer/model"
	"invoice-api/internal/features/customer/query"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerController struct {
	Command command.Command
	Query   query.Query
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
	keyword := c.Query("keyword")
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
	items, err := s.Query.GetItemsByQuery(keyword, size, page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
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

func (s *CustomerController) GetCustomersCount(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	keyword := c.Query("keyword")
	total_items, err := s.Query.GetTotalItemsByQuery(keyword)
	if err != nil {
		return c.JSON(0)
	}

	return c.JSON(total_items)
}

func (s *CustomerController) GetCustomersWithTotalByQuery(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	keyword := c.Query("keyword")
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

	res, err := s.Query.GetItemsWithTotalByQuery(keyword, size, page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(res)
}
