package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreatePayslipRequest represents the request body for creating a payslip
type CreatePayslipRequest struct {
	PeriodID string `json:"period_id" example:"0370f6f0-5762-4bbe-a57c-b398c2c76dd2"`
}

// CreatePayslip godoc
// @Summary Generate a new payslip
// @Description Employee can generate a payslip by providing period id
// @Tags Payslip
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param payslip body CreatePayslipRequest true "Payslip Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /payslip [post]
func (controller *Controller) CreatePayslip(c *fiber.Ctx) error {
	params := &service.CreatePayslipParam{}

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
	params.UserID = userID.(string)

	result, err := controller.service.CreatePayslip(context.Background(), params)
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

// GetPayslipSummary godoc
// @Summary Get payslip summary
// @Description Get payslip summary
// @Tags Payslip
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param start_date query string true "Start date" example:"2025-01-01"
// @Param end_date query string true "End date" example:"2025-01-31"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of attendance periods per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/payslip-summary [get]
func (controller *Controller) GetPayslipSummary(c *fiber.Ctx) error {
	params := &service.GetPayslipSummaryParam{}

	params.StartDate = c.Query("start_date")
	params.EndDate = c.Query("end_date")
	if page := c.Query("page"); page != "" {
		if pageNum, err := strconv.ParseInt(page, 10, 64); err == nil {
			params.Page = pageNum
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if pageSizeNum, err := strconv.ParseInt(pageSize, 10, 64); err == nil {
			params.PageSize = pageSizeNum
		}
	}

	result, err := controller.service.GetPayslipSummary(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"remark":         err.Error(),
			"status_code":    result.StatusCode,
			"status_message": result.StatusMessage,
		})
	}

	return c.JSON(result)
}
