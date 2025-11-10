package controller

import (
	"invoice-api/internal/features/revenue/command"
	"invoice-api/internal/features/revenue/model"
	"invoice-api/internal/features/revenue/query"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type RevenueController struct {
	Command           command.RevenueCommand
	Query             query.RevenueQuery
}

func (s *RevenueController) CreateRevenue(c *fiber.Ctx) error {
	s.Command = &command.DefaultRevenueCommand{}
	
	payload := new(model.CreateRevenue)
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
			"error": "Failed to create revenue",
		})
	}

	return c.Status(201).JSON(resp)
}

func (s *RevenueController) GetAllRevenues(c *fiber.Ctx) error {
	s.Query = &query.DefaultRevenueQuery{}
	items, err := s.Query.GetItemsByQuery()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch revenues",
		})
	}

	return c.JSON(items)
}

func (s *RevenueController) GetRevenueByID(c *fiber.Ctx) error {
	s.Query = &query.DefaultRevenueQuery{}
	id := c.Params("id")

	item, err := s.Query.GetItemByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Revenue not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch revenue",
		})
	}

	return c.JSON(item)
}

func (s *RevenueController) UpdateRevenue(c *fiber.Ctx) error {
	s.Command = &command.DefaultRevenueCommand{}
	id := c.Params("id")

	payload := new(model.UpdateRevenue)
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
				"error": "Revenue not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update revenue",
		})
	}

	if res.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to update revenue",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Revenue updated successfully",
	})
}

func (s *RevenueController) DeleteRevenue(c *fiber.Ctx) error {
	s.Command = &command.DefaultRevenueCommand{}
	id := c.Params("id")

	res, err := s.Command.DeleteItem(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Revenue not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update revenue",
		})
	}

	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Failed to delete revenue",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Revenue deleted successfully",
	})
}