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

type KontrolBalikController struct {
	KontrolBalikService *service.KontrolBalikService
	Modifier            *mold.Transformer
}

func NewKontrolBalikController(kontrolBalikService *service.KontrolBalikService, modifier *mold.Transformer) *KontrolBalikController {
	return &KontrolBalikController{kontrolBalikService, modifier}
}

func (c *KontrolBalikController) Search(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas && auth.Role != constant.RolePengguna {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikSearchRequest)
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	} else if auth.Role == constant.RolePengguna {
		request.IdPengguna = auth.ID
	}
	request.Status = ctx.Query("status")
	response, err := c.KontrolBalikService.Search(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}

func (c *KontrolBalikController) Get(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikGetRequest)
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
	response, err := c.KontrolBalikService.Get(ctx.Context(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}

func (c *KontrolBalikController) Create(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikCreateRequest)
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}
	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	if err := c.KontrolBalikService.Create(ctx.Context(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "Kontrol balik berhasil dibuat"})
}

func (c *KontrolBalikController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikUpdateRequest)
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

	if auth.Role == constant.RoleAdminPuskesmas {
		request.IdAdminPuskesmas = auth.ID
	}
	if err := c.KontrolBalikService.Update(ctx.Context(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Kontrol balik berhasil diperbarui"})
}

func (c *KontrolBalikController) Delete(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikDeleteRequest)
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
	if err := c.KontrolBalikService.Delete(ctx.Context(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Kontrol balik berhasil dihapus"})
}

func (c *KontrolBalikController) Batal(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikBatalRequest)
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
	if err := c.KontrolBalikService.Batal(ctx.Context(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Kontrol balik berhasil ditandai batal"})
}

func (c *KontrolBalikController) Selesai(ctx fiber.Ctx) error {
	auth := middleware.GetAuth(ctx)
	if auth.Role != constant.RoleAdminSuper && auth.Role != constant.RoleAdminPuskesmas {
		return fiber.ErrForbidden
	}
	request := new(model.KontrolBalikSelesaiRequest)
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
	if err := c.KontrolBalikService.Selesai(ctx.Context(), request); err != nil {
		log.Println(err.Error())
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Kontrol balik berhasil ditandai selesai"})
}
