package main

import (
	"payslip-generator/controller"

	"github.com/gofiber/fiber/v2"
)

func defineEndpoints(app *fiber.App, controller *controller.Controller) *fiber.App {
	app.Get("/ping", controller.Ping)

	return app
}