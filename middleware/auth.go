package middleware

import (
	"payslip-generator/util/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
	"payslip-generator/util/errs"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status_code":    errs.CODE_ERR_UNAUTHENTICATED,
				"status_message": "Missing Authorization Header",
				"remark":         "Missing Authorization Header",
			})
		}

		// Expects: Bearer <token>
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status_code":    errs.CODE_ERR_UNAUTHORIZED,
				"status_message": "Invalid Authorization header format",
				"remark":         "Invalid Authorization header format",
			})
		}

		token := splitToken[1]
		claims, err := auth.ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status_code":    errs.CODE_ERR_UNAUTHENTICATED,
				"status_message": "Invalid or expired token",
				"remark":         "Invalid or expired token",
			})
		}

		c.Locals("id", claims["user_id"])
		c.Locals("username", claims["username"])

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("username")
		if username != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status_code":    errs.CODE_ERR_UNAUTHORIZED,
				"status_message": "Admin only",
				"remark":         "Admin only",
			})
		}
		return c.Next()
	}
}
