package route

import (
	"invoice-api/internal/features/auth/controller"
	"invoice-api/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct{}

func (s *AuthRoute) Init(router *fiber.App) {
	authController := new(controller.AuthController)
	auth := router.Group("/auth")

	auth.Post("/signup", authController.SignUpUser)
	auth.Post("/signin", authController.SignInUser)
	auth.Get("/signout", middleware.Authorize, authController.LogoutUser)
}
