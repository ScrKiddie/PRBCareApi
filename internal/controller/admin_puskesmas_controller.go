package controller

import (
	"github.com/go-playground/mold/v4"
	"github.com/gofiber/fiber/v3"
	"log"
	"math"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
	"strconv"
)

type AdminPuskesmasController struct {
	AdminPuskesmasService *service.AdminPuskesmasService
	Modifier              *mold.Transformer
}

func NewAdminPuskesmasController(puskesmasService *service.AdminPuskesmasService, modifier *mold.Transformer) *AdminPuskesmasController {
	return &AdminPuskesmasController{puskesmasService, modifier}
}

func (c *AdminPuskesmasController) Login(ctx fiber.Ctx) error {
	request := new(model.AdminPuskesmasLoginRequest)
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	response, err := c.AdminPuskesmasService.Login(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": response.Token})
}

func (c *AdminPuskesmasController) Current(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasGetRequest)
	request.ID = auth.ID
	response, err := c.AdminPuskesmasService.Current(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminPuskesmasController) CurrentProfileUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasProfileUpdateRequest)
	request.ID = auth.ID
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if err := c.AdminPuskesmasService.CurrentProfileUpdate(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "Admin puskesmas berhasil diupdate"})
}

func (c *AdminPuskesmasController) CurrentPasswordUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasPasswordUpdateRequest)
	request.ID = auth.ID
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if err := c.AdminPuskesmasService.CurrentPasswordUpdate(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "Password berhasil diganti"})
}

func (c *AdminPuskesmasController) List(ctx fiber.Ctx) error {
	response, err := c.AdminPuskesmasService.List(ctx.Context())
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminPuskesmasController) Get(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasGetRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		log.Println("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request.ID = int32(id)
	response, err := c.AdminPuskesmasService.Get(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *AdminPuskesmasController) Create(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasCreateRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.AdminPuskesmasService.Create(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Admin puskesmas berhasil dibuat"})
}

func (c *AdminPuskesmasController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasUpdateRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		log.Println("value out of range for int32")
		return fiber.ErrBadRequest
	}

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	request.ID = int32(id)

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.AdminPuskesmasService.Update(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Admin puskesmas berhasil diupdate"})
}

func (c *AdminPuskesmasController) Delete(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper {
		return fiber.ErrForbidden
	}
	request := new(model.AdminPuskesmasDeleteRequest)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if id < math.MinInt32 || id > math.MaxInt32 {
		log.Println("value out of range for int32")
		return fiber.ErrBadRequest
	}
	request.ID = int32(id)

	if err := c.AdminPuskesmasService.Delete(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Admin puskesmas berhasil dihapus"})
}
