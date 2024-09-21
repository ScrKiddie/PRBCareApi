package controller

import (
	"github.com/go-playground/mold/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"log/slog"
	"math"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
	"strconv"
)

type ArtikelController struct {
	ArtikelService *service.ArtikelService
	Modifier       *mold.Transformer
}

func NewArtikelController(artikelService *service.ArtikelService, modifier *mold.Transformer) *ArtikelController {
	return &ArtikelController{
		ArtikelService: artikelService,
		Modifier:       modifier,
	}
}

func (c *ArtikelController) Get(ctx fiber.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		slog.Error("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request := new(model.ArtikelGetRequest)
	request.ID = int32(id)

	response, err := c.ArtikelService.Get(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (c *ArtikelController) Search(ctx fiber.Ctx) error {
	param := ctx.Query("idAdminPuskesmas")
	request := new(model.ArtikelSearchRequest)
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

	response, err := c.ArtikelService.Search(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (c *ArtikelController) Create(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}

	request := new(model.ArtikelCreateRequest)
	if err := ctx.Bind().Body(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}
	var file *model.File
	banner, err := ctx.FormFile("banner")
	if err != nil && err != fasthttp.ErrMissingFile {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	if banner != nil {
		file = &model.File{}
		file.FileHeader = banner
	}

	if err := c.ArtikelService.Create(ctx.Context(), request, file); err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Artikel berhasil dibuat"})
}

func (c *ArtikelController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
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

	request := new(model.ArtikelUpdateRequest)
	if err := ctx.Bind().Body(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	request.ID = int32(id)

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
		request.CurrentAdminPuskesmas = true
	}

	var file *model.File
	banner, err := ctx.FormFile("banner")
	if err != nil && err != fasthttp.ErrMissingFile {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	if banner != nil {
		file = &model.File{}
		file.FileHeader = banner
	}

	err = c.ArtikelService.Update(ctx.Context(), request, file)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Artikel berhasil diupdate"})
}

func (c *ArtikelController) Delete(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}

	request := new(model.ArtikelDeleteRequest)
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

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}

	err = c.ArtikelService.Delete(ctx.Context(), request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Artikel berhasil dihapus"})
}
