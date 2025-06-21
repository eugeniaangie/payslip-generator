package controller

import (
	"context"
	"payslip-generator/service"
	"payslip-generator/util/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username     string `json:"username" example:"john_doe"`
	PasswordHash string `json:"password_hash" example:"$2a$10$hashedpassword"`
	FullName     string `json:"full_name" example:"John Doe"`
	UserRole     string `json:"user_role" example:"employee"`
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags User
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User"
// @Success 200 {object} map[string]interface{}
// @Router /user [post]
func (controller *Controller) CreateUser(c *fiber.Ctx) error {
	params := &service.CreateUserParam{}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	result, err := controller.service.CreateUser(context.Background(), params)
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

// GetUserList godoc
// @Summary Get user list
// @Description Get user list
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Number of users per page"
// @Success 200 {object} map[string]interface{}
// @Router /user [get]
func (controller *Controller) GetUserList(c *fiber.Ctx) error {
	params := &service.GetUserListParam{}

	params.Search = c.Query("search")

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

	result, err := controller.service.GetUserList(context.Background(), params)
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
