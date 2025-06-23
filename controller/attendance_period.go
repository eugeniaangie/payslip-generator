package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateAttendancePeriodRequest represents the request body for creating a attendance period
type CreateAttendancePeriodRequest struct {
	StartDate string `json:"start_date" example:"2025-01-01"`
	EndDate   string `json:"end_date" example:"2025-01-31"`
}

// CreateAttendancePeriod godoc
// @Summary Create a new attendance period
// @Description Only admin can create an attendance period by providing start and end dates
// @Tags Attendance Period
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param attendance_period body CreateAttendancePeriodRequest true "Attendance Period Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/attendance_period [post]
func (controller *Controller) CreateAttendancePeriod(c *fiber.Ctx) error {
	params := &service.CreateAttendancePeriodParam{}

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

	result, err := controller.service.CreateAttendancePeriod(context.Background(), params)
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

// GetAttendancePeriodList godoc
// @Summary Get attendance period list
// @Description Get attendance period list
// @Tags Attendance Period
// @Accept json
// @Produce json
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of attendance periods per page"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /attendance_period [get]
func (controller *Controller) GetAttendancePeriodList(c *fiber.Ctx) error {
	params := &service.GetAttendancePeriodListParam{}

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

	result, err := controller.service.GetAttendancePeriodList(context.Background(), params)
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
		"pagination":     result.Pagination,
	}

	return c.JSON(response)
}
