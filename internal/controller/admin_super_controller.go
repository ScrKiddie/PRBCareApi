package controller

import (
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
)

type AdminSuperController struct {
	AdminSuperService *service.AdminSuperService
}

func NewAdminSuperController(adminSuperService *service.AdminSuperService) *AdminSuperController {
	return &AdminSuperController{adminSuperService}
}

func (c *AdminSuperController) Login(ctx fiber.Ctx) error {
	request := new(model.AdminSuperLoginRequest)
	if err := ctx.Bind().JSON(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	response, err := c.AdminSuperService.Login(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": response.Token})
}

func (c *AdminSuperController) PasswordUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper {
		return fiber.ErrForbidden
	}
	request := new(model.AdminSuperPasswordUpdateRequest)
	request.ID = auth.ID
	if err := ctx.Bind().JSON(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if err := c.AdminSuperService.PasswordUpdate(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "Password berhasil diganti"})
}
