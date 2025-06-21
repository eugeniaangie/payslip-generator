package controller

import (
	"context"

	"payslip-generator/service"

	"github.com/gofiber/fiber/v2"
)

// Ping godoc
// @Summary Health check
// @Description Returns a basic pong message
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func (controller *Controller) Ping(c *fiber.Ctx) error {
	params := &service.PingParams{}

	// call service layer
	result, err := controller.service.Ping(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"remark": err.Error(),
		})
	}

	// tidy up response
	response := map[string]interface{}{
		"status_code":    result.StatusCode,
		"status_message": result.StatusMessage,
		"data": map[string]interface{}{
			"message": result.Data.Message,
		},
	}

	return c.JSON(response)
}
