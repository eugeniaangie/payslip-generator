package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

type CreateSalaryRequest struct {
	UserID   string `json:"user_id" example:"369e9c16-2df2-4a18-ab83-d6c0065ecfce"`
	Amount   int    `json:"amount" example:"100000"`
	IsActive bool   `json:"is_active" example:"true"`
}

// CreateSalary godoc
// @Summary Create a new salary
// @Description Only admin can create a salary by providing user id, amount, is active, ip address and created by
// @Tags Salary
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateSalaryRequest true "Salary request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/salary [post]
func (controller *Controller) CreateSalary(c *fiber.Ctx) error {
	params := &service.CreateSalaryParam{}

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

	result, err := controller.service.CreateSalary(context.Background(), params)
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
