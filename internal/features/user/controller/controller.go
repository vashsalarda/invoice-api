package controller

import (
	"invoice-api/internal/features/user/command"
	"invoice-api/internal/features/user/model"
	"invoice-api/internal/features/user/query"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	Command           command.Command
	Query             query.Query
}

// HandleCreateUser handles the HTTP request to create a user
func (s *UserController) CreateUser(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	s.Command = &command.DefaultCommand{}
	
	payload := new(model.CreateUser)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}

	item, _ := s.Query.GetByEmail(payload.Email)
	if item.ID != primitive.NilObjectID {
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "The email address already taken. Please select another email address"})
	}

	resp, err := s.Command.CreateUser(payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(resp)
}

// HandleGetAllUsers handles the HTTP request to get all users
func (s *UserController) GetAllUsers(c *fiber.Ctx) error {
	s.Query = &query.DefaultQuery{}
	users, err := s.Query.GetAllByQuery()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// HandleGetUserByID handles the HTTP request to get a user by ID
func (s *UserController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.Query.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

// HandleUpdateUser handles the HTTP request to update a user
func (s *UserController) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	payload := new(model.UpdateUser)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := model.ValidateStruct(payload)
	if validationErrors != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "errors": validationErrors})
	}
	res, err := s.Command.UpdateUser(id, payload)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
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
		"message": "User updated successfully",
	})
}

// HandleDeleteUser handles the HTTP request to delete a user
func (s *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := s.Command.DeleteUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
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
		"message": "User deleted successfully",
	})
}