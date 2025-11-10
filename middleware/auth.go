package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"invoice-api/internal/database"
	"invoice-api/internal/features/user/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func Authorize(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if after, ok := strings.CutPrefix(authorization, "Bearer "); ok {
		tokenString = after
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	if jwtSecret == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Warning: Failed to load .env file")
		}
		jwtSecret = os.Getenv("JWT_SECRET")
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("Invalid token [01]: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid token claim"})
	}

	var user model.User
	db := database.GetDatabase()
	collection := db.Collection("users")

	var claimsSub = fmt.Sprint(claims["sub"])
	var error error
	ID, _ := primitive.ObjectIDFromHex(claimsSub)
	error = collection.FindOne(context.TODO(), bson.M{"_id": ID}).Decode(&user)

	if error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("Invalid token [02]: %v", error)})
	}

	if user.ID.Hex() != claimsSub {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "The user belonging to this token no logger exists"})
	}

	c.Locals("user", model.FilterUserRecord(&user))

	return c.Next()
}
