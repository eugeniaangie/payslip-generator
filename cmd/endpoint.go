package main

import (
	"payslip-generator/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
)

func defineEndpoints(app *fiber.App, controller *controller.Controller) *fiber.App {
	// Add Swagger handler
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")
	api.Get("/ping", controller.Ping)
	api.Post("/user", controller.CreateUser)
	api.Get("/user", controller.GetUserList)

	return app
}
