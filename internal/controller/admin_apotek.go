package controller

import (
	"github.com/go-playground/mold/v4"
	"github.com/gofiber/fiber/v3"
	"log"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
	"strconv"
)

type AdminApotekController struct {
	AdminApotekService *service.AdminApotekService
	Modifier           *mold.Transformer
}

func NewAdminApotekController(apotekService *service.AdminApotekService, modifier *mold.Transformer) *AdminApotekController {
	return &AdminApotekController{apotekService, modifier}
}

func (c *AdminApotekController) Login(ctx fiber.Ctx) error {
	request := new(model.AdminApotekLoginRequest)
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	response, err := c.AdminApotekService.Login(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": response.Token})
}

func (c *AdminApotekController) Current(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	request := new(model.AdminApotekGetRequest)
	request.ID = auth.ID
	response, err := c.AdminApotekService.Current(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminApotekController) CurrentProfileUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	request := new(model.AdminApotekProfileUpdateRequest)
	request.ID = auth.ID
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if err := c.AdminApotekService.CurrentProfileUpdate(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "Admin apotek berhasil diupdate"})
}

func (c *AdminApotekController) CurrentPasswordUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	request := new(model.AdminApotekPasswordUpdateRequest)
	request.ID = auth.ID
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if err := c.AdminApotekService.CurrentPasswordUpdate(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "Password berhasil diganti"})
}

func (c *AdminApotekController) List(ctx fiber.Ctx) error {
	response, err := c.AdminApotekService.List(ctx.Context())
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminApotekController) Get(ctx fiber.Ctx) error {
	request := new(model.AdminApotekGetRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	request.ID = id
	response, err := c.AdminApotekService.Get(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminApotekController) Create(ctx fiber.Ctx) error {
	request := new(model.AdminApotekCreateRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.AdminApotekService.Create(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Admin apotek berhasil dibuat"})
}

func (c *AdminApotekController) Update(ctx fiber.Ctx) error {
	request := new(model.AdminApotekUpdateRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	request.ID = id
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.AdminApotekService.Update(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Admin apotek berhasil diupdate"})
}

func (c *AdminApotekController) Delete(ctx fiber.Ctx) error {
	request := new(model.AdminApotekDeleteRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	request.ID = id

	if err := c.AdminApotekService.Delete(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Admin apotek berhasil dihapus"})
}
