package controller

import (
	"fmt"
	"os"
	"time"

	"invoice-api/internal/features/user/command"
	"invoice-api/internal/features/user/model"
	"invoice-api/internal/features/user/query"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Command command.Command
	Query   query.Query
}

func (s *AuthController) SignUpUser(c *fiber.Ctx) error {
	var payload *model.CreateUser
	if s.Query == nil {
		s.Query = &query.DefaultQuery{}
	}
	if s.Command == nil {
		s.Command = &command.DefaultCommand{}
	}

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("error: %+v\n", err)
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := model.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Some fields are required. Please fill in the required fields", "errors": errors})
	}

	user, err := s.Query.GetByEmail(payload.Email)
	if user != (model.User{}) {
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "The email address already exists. Please select another email address"})
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	_, err = s.Command.CreateUser(payload)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "User created successfully"})
}

func (s *AuthController) SignInUser(c *fiber.Ctx) error {
	var payload *model.SignIn
	if s.Query == nil {
		s.Query = &query.DefaultQuery{}
	}

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("Error: %+v\n", err)
		return c.Status(422).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := model.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Some fields are required. Please fill in the required fields", "errors": errors})
	}

	user, err := s.Query.GetByEmail(payload.Email)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["exp"] = now.Add(time.Hour * 24).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   60 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func (s *AuthController) LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// func (s *AuthController) GetUsers(c *fiber.Ctx) error {
// 	if s.Query == nil {
// 		s.Query = &query.DefaultQuery{}
// 	}
// 	keyword := c.Query("keyword")
// 	kind := c.Query("kind")
// 	sizeStr := c.Query("size")
// 	pageStr := c.Query("page")
// 	size, err := strconv.ParseInt(sizeStr, 10, 64)
// 	if err != nil {
// 		size = 25
// 	}
// 	page, err := strconv.ParseInt(pageStr, 10, 64)
// 	if err != nil {
// 		page = 1
// 	}
// 	totalItems, err := s.Query.CountDocuments(keyword, kind)
// 	if err != nil {
// 		return c.Status(400).SendString(err.Error())
// 	}
// 	items, err := s.Query.FindAll(keyword, kind, size, page)
// 	if err != nil {
// 		return c.Status(400).SendString(err.Error())
// 	}

// 	resp := model.UserPage{
// 		PageNumber: page,
// 		PageSize:   size,
// 		TotalRows:  totalItems,
// 		Docs:       items,
// 	}

// 	return c.Status(200).JSON(resp)
// }

func (s *AuthController) GetUser(c *fiber.Ctx) error {
	if s.Query == nil {
		s.Query = &query.DefaultQuery{}
	}

	id := c.Params("id")

	user, err := s.Query.GetByID(id)

	if user.ID.String() == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.Status(200).JSON(user)
}
