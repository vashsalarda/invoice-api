package server

import (
	"log"

	"github.com/gofiber/fiber/v2"

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
