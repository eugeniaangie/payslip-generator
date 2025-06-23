package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

// CreateAttendance godoc
// @Summary Create a new attendance
// @Description Only admin can create an attendance by providing user id, date, ip address and created by
// @Tags Attendance
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /attendance [post]
func (controller *Controller) CreateAttendance(c *fiber.Ctx) error {
	params := &service.CreateAttendanceParam{}

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

	result, err := controller.service.CreateAttendance(context.Background(), params)
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
