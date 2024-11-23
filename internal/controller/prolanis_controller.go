package controller

import (
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"math"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
	"strconv"
)

type ProlanisController struct {
	ProlanisService *service.ProlanisService
}

func NewProlanisController(prolanisService *service.ProlanisService) *ProlanisController {
	return &ProlanisController{prolanisService}
}

func (c *ProlanisController) Search(ctx fiber.Ctx) error {
	request := new(model.ProlanisSearchRequest)
	param := ctx.Query("idAdminPuskesmas")
	if param != "" {
		idAdminPuskesmas, err := strconv.Atoi(param)
		if err != nil {
			slog.Error(err.Error())
			return fiber.ErrBadRequest
		}
		if idAdminPuskesmas < math.MinInt32 || idAdminPuskesmas > math.MaxInt32 {
			slog.Error("value out of range for int32")
			return fiber.ErrBadRequest
		}
		request.IdAdminPuskesmas = int32(idAdminPuskesmas)
	}
	request.Status = ctx.Query("status")
	response, err := c.ProlanisService.Search(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{

		"data": response})
}

func (c *ProlanisController) Get(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.ProlanisGetRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		slog.Error("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request.ID = int32(id)
	response, err := c.ProlanisService.Get(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *ProlanisController) Create(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.ProlanisCreateRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	if err := c.ProlanisService.Create(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Prolanis berhasil dibuat"})
}

func (c *ProlanisController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.ProlanisUpdateRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		slog.Error("value out of range for int32")
		return fiber.ErrBadRequest
	}

	if err := ctx.Bind().JSON(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	request.ID = int32(id)

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
		request.CurrentAdminPuskesmas = true
	}

	if err := c.ProlanisService.Update(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Prolanis berhasil diupdate"})
}

func (c *ProlanisController) Delete(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.ProlanisDeleteRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		slog.Error("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request.ID = int32(id)

	if err := c.ProlanisService.Delete(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Prolanis berhasil dihapus"})
}

func (c *ProlanisController) Selesai(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.ProlanisSelesaiRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		slog.Error("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request.ID = int32(id)

	if err := c.ProlanisService.Selesai(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Prolanis berhasil ditandai selesai"})
}
