package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"

	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	UserRole     string `json:"user_role"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User"
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func (controller *Controller) Register(c *fiber.Ctx) error {
	params := &service.RegisterParam{}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	result, err := controller.service.Register(context.Background(), params)
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "User"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func (controller *Controller) Login(c *fiber.Ctx) error {
	params := &service.LoginParam{}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	result, err := controller.service.Login(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"remark":         err.Error(),
			"status_code":    errs.CODE_ERR_DATABASE,
			"status_message": "failed to login",
		})
	}

	response := map[string]interface{}{
		"status_code":    errs.CODE_SUCCESS,
		"status_message": "ok",
		"data":           result,
	}

	return c.JSON(response)
}