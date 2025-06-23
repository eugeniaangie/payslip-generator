package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

// RunPayrollRequest represents the request body for running payroll
type RunPayrollRequest struct {
	PeriodID  string `json:"period_id" example:"0370f6f0-5762-4bbe-a57c-b398c2c76dd2"`
}

// RunPayroll godoc
// @Summary Run payroll
// @Description Admin can run payroll by providing period id
// @Tags Payroll
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param payroll body RunPayrollRequest true "Payroll Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/payroll [post]
func (controller *Controller) RunPayroll(c *fiber.Ctx) error {
	params := &service.RunPayrollParam{}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	userID := c.Locals("id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status_code":    errs.CODE_ERR_UNAUTHORIZED,
			"status_message": "unauthorized",
			"remark":         "unauthorized",
		})
	}

	params.CreatedBy = userID.(string)
	params.UpdatedBy = userID.(string)
	params.IpAddress = c.IP()

	result, err := controller.service.RunPayroll(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"remark":         err.Error(),
			"status_code":    result.StatusCode,
			"status_message": result.StatusMessage,
		})
	}

	response := map[string]interface{}{
		"status_code":    result.StatusCode,
		"status_message": result.StatusMessage,
		"data":           result.Data,
	}

	return c.JSON(response)
}
