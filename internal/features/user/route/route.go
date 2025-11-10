package route

import (
	"invoice-api/internal/features/user/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type UserRoute struct{}

func (c *UserRoute) Init(router *fiber.App) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	router.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Singapore",
		Format:     "${blue}[${time}] | ${green}${status} | ${cyan}${latency} | ${blue}${ip} | ${method} | ${white}${path} | ${red}${error}${white}\n",
	}))
	
	controller := new(controller.UserController)
	users := router.Group("/users")

	users.Post("/", controller.CreateUser)
	users.Get("/", controller.GetAllUsers)
	users.Get("/:id", controller.GetUserByID)
	users.Patch("/:id", controller.UpdateUser)
	users.Delete("/:id", controller.DeleteUser)
}