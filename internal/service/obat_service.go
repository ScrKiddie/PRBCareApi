package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
)

type ObatService struct {
	DB                        *gorm.DB
	ObatRepository            *repository.ObatRepository
	AdminApotekRepository     *repository.AdminApotekRepository
	PengambilanObatRepository *repository.PengambilanObatRepository
	Validator                 *validator.Validate
}

func NewObatService(
	db *gorm.DB,
	obatRepository *repository.ObatRepository,
	adminApotekRepository *repository.AdminApotekRepository,
	pengambilanObatRepository *repository.PengambilanObatRepository,
	validator *validator.Validate) *ObatService {
	return &ObatService{db, obatRepository, adminApotekRepository, pengambilanObatRepository, validator}
}

func (s *ObatService) List(ctx context.Context, request *model.ObatListRequest) (*[]model.ObatResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	obat := new([]entity.Obat)

	if request.IdAdminApotek > 0 {
		if err := s.ObatRepository.FindAllByIdAdminApotek(tx, obat, request.IdAdminApotek); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if err := s.ObatRepository.FindAll(tx, obat); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	var response []model.ObatResponse
	for _, o := range *obat {
		response = append(response, model.ObatResponse{
			ID:       o.ID,
			NamaObat: o.NamaObat,
			Jumlah:   o.Jumlah,
			AdminApotek: &model.AdminApotekResponse{
				ID:         o.AdminApotek.ID,
				NamaApotek: o.AdminApotek.NamaApotek,
				Telepon:    o.AdminApotek.Telepon,
				Alamat:     o.AdminApotek.Alamat,
			},
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *ObatService) Get(ctx context.Context, request *model.ObatGetRequest) (*model.ObatResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	obat := new(entity.Obat)
	if request.IdAdminApotek > 0 {
		if err := s.ObatRepository.FindByIdAndIdAdminApotek(tx, obat, request.ID, request.IdAdminApotek); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.ObatRepository.FindById(tx, obat, request.ID); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.ObatResponse)
	response.ID = obat.ID
	response.NamaObat = obat.NamaObat
	response.Jumlah = obat.Jumlah
	response.IdAdminApotek = obat.IdAdminApotek

	return response, nil
}

func (s *ObatService) Create(ctx context.Context, request *model.ObatCreateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := s.AdminApotekRepository.FindById(tx, &entity.AdminApotek{}, request.IdAdminApotek); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	obatEnity := new(entity.Obat)
	obatEnity.NamaObat = request.NamaObat
	obatEnity.IdAdminApotek = request.IdAdminApotek
	obatEnity.Jumlah = request.Jumlah

	if err := s.ObatRepository.Create(tx, obatEnity); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ObatService) Update(ctx context.Context, request *model.ObatUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := s.AdminApotekRepository.FindById(tx, &entity.AdminApotek{}, request.IdAdminApotek); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	obat := new(entity.Obat)
	if request.CurrentAdminApotek {
		if err := s.ObatRepository.FindByIdAndIdAdminApotekAndLockForUpdate(tx, obat, request.ID, request.IdAdminApotek); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.ObatRepository.FindByIdAndLockForUpdate(tx, obat, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	obat.IdAdminApotek = request.IdAdminApotek
	obat.NamaObat = request.NamaObat
	obat.Jumlah = request.Jumlah

	if err := s.ObatRepository.Update(tx, obat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ObatService) Delete(ctx context.Context, request *model.ObatDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	obat := new(entity.Obat)
	if request.IdAdminApotek > 0 {
		if err := s.ObatRepository.FindByIdAndIdAdminApotek(tx, obat, request.ID, request.IdAdminApotek); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.ObatRepository.FindById(tx, obat, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.PengambilanObatRepository.FindByIdObat(tx, &entity.PengambilanObat{}, request.ID); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Obat masih terkait dengan data pengambilan obat yang ada")
	}

	if err := s.ObatRepository.Delete(tx, obat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
