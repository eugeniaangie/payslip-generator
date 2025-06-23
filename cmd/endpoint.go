package main

import (
	"payslip-generator/controller"
	"payslip-generator/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
)

func defineEndpoints(app *fiber.App, controller *controller.Controller) *fiber.App {
	// Add Swagger handler
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")
	api.Get("/ping", controller.Ping)

	auth := api.Group("/auth")
	auth.Post("/login", controller.Login)
	auth.Post("/register", controller.Register)

	apiAdmin := api.Group("/admin")
	apiAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	apiAdmin.Post("/salary", controller.CreateSalary)
	apiAdmin.Post("/attendance_period", controller.CreateAttendancePeriod)
	// protected.Post("/run-payroll", controller.RunPayroll)
	// protected.Get("/payslip-summary", controller.GetPayslipSummary)

	apiEmployee := api.Group("/")
	apiEmployee.Use(middleware.AuthMiddleware())
	apiEmployee.Post("/attendance", controller.CreateAttendance)
	apiEmployee.Post("/overtime", controller.CreateOvertime)
	apiEmployee.Post("/reimbursement", controller.CreateReimbursement)
	apiEmployee.Get("/attendance_period", controller.GetAttendancePeriodList)
	apiEmployee.Post("/payslip", controller.CreatePayslip)

	return app
}
