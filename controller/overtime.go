package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

// CreateOvertimeRequest represents the request body for creating a overtime
type CreateOvertimeRequest struct {
	Hours int `json:"hours" example:"1"`
}

// CreateOvertime godoc
// @Summary Create a new overtime
// @Description Only admin can create an overtime by providing user id, hours, ip address and created by
// @Tags Overtime
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param overtime body CreateOvertimeRequest true "Overtime Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /overtime [post]
func (controller *Controller) CreateOvertime(c *fiber.Ctx) error {
	params := &service.CreateOvertimeParam{}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	userID := c.Locals("id")
	ipAddress := c.IP()

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status_code":    errs.CODE_ERR_UNAUTHORIZED,
			"status_message": "unauthorized",
			"remark":         "unauthorized",
		})
	}

	params.UserID = userID.(string)
	params.IpAddress = ipAddress
	params.CreatedBy = userID.(string)
	params.UpdatedBy = userID.(string)

	result, err := controller.service.CreateOvertime(context.Background(), params)
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
