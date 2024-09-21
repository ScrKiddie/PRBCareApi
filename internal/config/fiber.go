package config

import (
	"github.com/gofiber/fiber/v3"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: ErrorHandler(), BodyLimit: 50 * 1024 * 1024})
}
func ErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {
		code := err.(*fiber.Error).Code
		if code == fiber.StatusNotFound {
			return ctx.Status(code).JSON(fiber.Map{
				"error": "Not found",
			})
		}
		return ctx.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
