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
	apiAdmin.Get("/attendance_period", controller.GetAttendancePeriodList)
	// protected.Post("/run-payroll", controller.RunPayroll)
	// protected.Get("/payslip-summary", controller.GetPayslipSummary)

	apiEmployee := api.Group("/")
	apiEmployee.Use(middleware.AuthMiddleware())
	apiEmployee.Post("/attendance", controller.CreateAttendance)
	apiEmployee.Post("/overtime", controller.CreateOvertime)
	apiEmployee.Post("/reimbursement", controller.CreateReimbursement)
	// apiEmployee.Get("/payslip/:period_id", controller.GetPayslip)

	// api.Post("/user", controller.CreateUser)
	// api.Get("/user", controller.GetUserList)
	// api.Get("/user/username", controller.GetUserByUsername)
	// api.Post("/attendance_period", controller.CreateAttendancePeriod)
	// api.Get("/attendance_period", controller.GetAttendancePeriodList)

	return app
}

// ADMIN ENDPOINTS (🔒 Admin-only)
// Method	Path	Description
// POST	/admin/attendance-period	Buat attendance period baru (start & end date)
// POST	/admin/run-payroll	Jalankan payroll untuk satu attendance period
// GET	/admin/payslip-summary	Lihat semua ringkasan payslip semua employee

// ✅ EMPLOYEE ENDPOINTS (👤 Authenticated Employee)
// Method	Path	Description
// POST	/attendance	Submit kehadiran
// POST	/overtime	Submit lembur (setelah jam kerja, max 3 jam/hari)
// POST	/reimbursement	Submit reimburse (dengan nominal & deskripsi)
// GET	/payslip/:period_id	Generate detail payslip untuk periode tertentu
