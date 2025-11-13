package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"invoice-api/internal/database"
	auth_route "invoice-api/internal/features/auth/route"
	customer_route "invoice-api/internal/features/customer/route"
	invoice_route "invoice-api/internal/features/invoice/route"
	revenue_route "invoice-api/internal/features/revenue/route"
	user_route "invoice-api/internal/features/user/route"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "invoice-api",
			AppName:      "invoice-api",
		}),

		db: database.New(),
	}

	server.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	server.App.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Singapore",
		Format:     "${blue}[${time}] | ${green}${status} | ${cyan}${latency} | ${blue}${ip} | ${method} | ${white}${path} | ${red}${error}${white}\n",
	}))

	userRoute := new(user_route.UserRoute)
	userRoute.Init(server.App)
	customerRoute := new(customer_route.CustomerRoute)
	customerRoute.Init(server.App)
	invoiceRoute := new(invoice_route.InvoiceRoute)
	invoiceRoute.Init(server.App)
	revenueRoute := new(revenue_route.InvoiceRoute)
	revenueRoute.Init(server.App)
	authRoute := new(auth_route.AuthRoute)
	authRoute.Init(server.App)

	// List all routes
	log.Println("API Routes:")
	routes := server.App.GetRoutes()
	for _, route := range routes {
		method := route.Method
		path := route.Path
		if method != "HEAD" {
			log.Printf("%-6s %-20s\n", method, path)
		}
	}

	return server
}
