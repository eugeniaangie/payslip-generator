package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

type CreateReimbursementRequest struct {
	Date        string `json:"date" example:"2025-01-01"`
	Amount      int    `json:"amount" example:"100000"`
	Description string `json:"description" example:"Transportation"`
}

// CreateReimbursement godoc
// @Summary Create a new reimbursement
// @Description Only admin can create a reimbursement by providing date, amount, description, ip address and created by
// @Tags Reimbursement
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateReimbursementRequest true "Reimbursement request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /reimbursement [post]
func (controller *Controller) CreateReimbursement(c *fiber.Ctx) error {
	params := &service.CreateReimbursementParam{}

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

	result, err := controller.service.CreateReimbursement(context.Background(), params)
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
