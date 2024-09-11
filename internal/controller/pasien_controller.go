package controller

import (
	"log"
	"math"
	"strconv"

	"github.com/go-playground/mold/v4"
	"github.com/gofiber/fiber/v3"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
)

type PasienController struct {
	PasienService *service.PasienService
	Modifier      *mold.Transformer
}

func NewPasienController(pasienService *service.PasienService, modifier *mold.Transformer) *PasienController {
	return &PasienController{pasienService, modifier}
}

func (c *PasienController) Search(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas && auth.Role != constant.RolePengguna {
		return fiber.ErrForbidden
	}
	request := new(model.PasienSearchRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	} else if auth.Role == constant.RolePengguna {
		request.IdPengguna = auth.ID
	}
	request.Status = ctx.Query("status")
	response, err := c.PasienService.Search(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *PasienController) Get(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.PasienGetRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
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
	response, err := c.PasienService.Get(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response})
}

func (c *PasienController) Create(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.PasienCreateRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.PasienService.Create(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Pasien berhasil dibuat"})
}

func (c *PasienController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.PasienUpdateRequest)
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

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
		request.CurrentAdminPuskesmas = true
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.PasienService.Update(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Pasien berhasil diupdate"})
}

func (c *PasienController) Delete(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.PasienDeleteRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
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

	if err := c.PasienService.Delete(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Pasien berhasil dihapus"})
}

func (c *PasienController) Selesai(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.PasienSelesaiRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
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

	if err := c.PasienService.Selesai(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Pasien berhasil ditandai selesai"})
}
