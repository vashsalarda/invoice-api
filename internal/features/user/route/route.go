package route

import (
	"invoice-api/internal/features/user/controller"

	"github.com/gofiber/fiber/v2"
)

type UserRoute struct{}

func (c *UserRoute) Init(router *fiber.App) {
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:3001",
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	// 	AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	// 	AllowCredentials: true,
	// }))
	
	controller := new(controller.UserController)
	users := router.Group("/users")

	users.Post("/", controller.CreateUser)
	users.Get("/", controller.GetAllUsers)
	users.Get("/:id", controller.GetUserByID)
	users.Patch("/:id", controller.UpdateUser)
	users.Delete("/:id", controller.DeleteUser)
}