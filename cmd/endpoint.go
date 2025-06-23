package main

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"payslip-generator/controller"
	"payslip-generator/middleware"
	"payslip-generator/service"
)

func defineEndpoints(app *fiber.App, controller *controller.Controller, service *service.Service) *fiber.App {
	// Add Swagger handler
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")
	api.Get("/ping", controller.Ping)

	auth := api.Group("/auth")
	auth.Post("/login", controller.Login)
	auth.Post("/register", controller.Register)

	apiAdmin := api.Group("/admin")
	apiAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(service))
	apiAdmin.Post("/salary", controller.CreateSalary)
	apiAdmin.Post("/attendance_period", controller.CreateAttendancePeriod)
	apiAdmin.Post("/payroll", controller.RunPayroll)
	apiAdmin.Get("/payslip-summary", controller.GetPayslipSummary)

	apiEmployee := api.Group("/")
	apiEmployee.Use(middleware.AuthMiddleware())
	apiEmployee.Post("/attendance", controller.CreateAttendance)
	apiEmployee.Post("/overtime", controller.CreateOvertime)
	apiEmployee.Post("/reimbursement", controller.CreateReimbursement)
	apiEmployee.Get("/attendance_period", controller.GetAttendancePeriodList)
	apiEmployee.Post("/payslip", controller.CreatePayslip)

	return app
}
