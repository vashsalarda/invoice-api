package controller

import (
	"invoice-api/internal/features/user/command"
	"invoice-api/internal/features/user/model"
	"invoice-api/internal/features/user/query"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	Command           command.Command
	Query             query.Query
}

// HandleCreateUser handles the HTTP request to create a user
func (s *UserController) CreateUser(c *fiber.Ctx) error {
	s.Command = &command.DefaultCommand{}
	var user *model.CreateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := s.Command.CreateUser(user)
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

	var val *model.UpdateUser
	if err := c.BodyParser(&val); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := s.Command.UpdateUser(id, val)
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